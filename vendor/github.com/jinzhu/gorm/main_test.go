package gorm_test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	testdb "github.com/erikstmartin/go-testdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"os"
	"testing"
	"time"
)

var (
	DB                 gorm.DB
	t1, t2, t3, t4, t5 time.Time
)

func init() {
	var err error
	switch os.Getenv("GORM_DIALECT") {
	case "mysql":
		// CREATE USER 'gorm'@'localhost' IDENTIFIED BY 'gorm';
		// CREATE DATABASE gorm;
		// GRANT ALL ON gorm.* TO 'gorm'@'localhost';
		fmt.Println("testing mysql...")
		DB, err = gorm.Open("mysql", "gorm:gorm@/gorm?charset=utf8&parseTime=True")
	case "postgres":
		fmt.Println("testing postgres...")
		DB, err = gorm.Open("postgres", "user=gorm DB.name=gorm sslmode=disable")
	case "foundation":
		fmt.Println("testing foundation...")
		DB, err = gorm.Open("foundation", "dbname=gorm port=15432 sslmode=disable")
	case "mssql":
		fmt.Println("testing mssql...")
		DB, err = gorm.Open("mssql", "server=SERVER_HERE;database=rogue;user id=USER_HERE;password=PW_HERE;port=1433")
	default:
		fmt.Println("testing sqlite3...")
		DB, err = gorm.Open("sqlite3", "/tmp/gorm.db")
	}

	// DB.SetLogger(Logger{log.New(os.Stdout, "\r\n", 0)})
	// DB.SetLogger(log.New(os.Stdout, "\r\n", 0))
	DB.LogMode(true)
	DB.LogMode(false)

	if err != nil {
		panic(fmt.Sprintf("No error should happen when connect database, but got %+v", err))
	}

	DB.DB().SetMaxIdleConns(10)

	runMigration()
}

func TestStringPrimaryKey(t *testing.T) {
	type UUIDStruct struct {
		ID   string `gorm:"primary_key"`
		Name string
	}
	DB.AutoMigrate(&UUIDStruct{})

	data := UUIDStruct{ID: "uuid", Name: "hello"}
	if err := DB.Save(&data).Error; err != nil || data.ID != "uuid" {
		t.Errorf("string primary key should not be populated")
	}
}

func TestExceptionsWithInvalidSql(t *testing.T) {
	var columns []string
	if DB.Where("sdsd.zaaa = ?", "sd;;;aa").Pluck("aaa", &columns).Error == nil {
		t.Errorf("Should got error with invalid SQL")
	}

	if DB.Model(&User{}).Where("sdsd.zaaa = ?", "sd;;;aa").Pluck("aaa", &columns).Error == nil {
		t.Errorf("Should got error with invalid SQL")
	}

	if DB.Where("sdsd.zaaa = ?", "sd;;;aa").Find(&User{}).Error == nil {
		t.Errorf("Should got error with invalid SQL")
	}

	var count1, count2 int64
	DB.Model(&User{}).Count(&count1)
	if count1 <= 0 {
		t.Errorf("Should find some users")
	}

	if DB.Where("name = ?", "jinzhu; delete * from users").First(&User{}).Error == nil {
		t.Errorf("Should got error with invalid SQL")
	}

	DB.Model(&User{}).Count(&count2)
	if count1 != count2 {
		t.Errorf("No user should not be deleted by invalid SQL")
	}
}

func TestSetTable(t *testing.T) {
	DB.Create(getPreparedUser("pluck_user1", "pluck_user"))
	DB.Create(getPreparedUser("pluck_user2", "pluck_user"))
	DB.Create(getPreparedUser("pluck_user3", "pluck_user"))

	if err := DB.Table("users").Where("role = ?", "pluck_user").Pluck("age", &[]int{}).Error; err != nil {
		t.Errorf("No errors should happen if set table for pluck", err.Error())
	}

	var users []User
	if DB.Table("users").Find(&[]User{}).Error != nil {
		t.Errorf("No errors should happen if set table for find")
	}

	if DB.Table("invalid_table").Find(&users).Error == nil {
		t.Errorf("Should got error when table is set to an invalid table")
	}

	DB.Exec("drop table deleted_users;")
	if DB.Table("deleted_users").CreateTable(&User{}).Error != nil {
		t.Errorf("Create table with specified table")
	}

	DB.Table("deleted_users").Save(&User{Name: "DeletedUser"})

	var deletedUsers []User
	DB.Table("deleted_users").Find(&deletedUsers)
	if len(deletedUsers) != 1 {
		t.Errorf("Query from specified table")
	}

	DB.Save(getPreparedUser("normal_user", "reset_table"))
	DB.Table("deleted_users").Save(getPreparedUser("deleted_user", "reset_table"))
	var user1, user2, user3 User
	DB.Where("role = ?", "reset_table").First(&user1).Table("deleted_users").First(&user2).Table("").First(&user3)
	if (user1.Name != "normal_user") || (user2.Name != "deleted_user") || (user3.Name != "normal_user") {
		t.Errorf("unset specified table with blank string")
	}
}

type Order struct {
}

type Cart struct {
}

func (c Cart) TableName() string {
	return "shopping_cart"
}

func TestHasTable(t *testing.T) {
	type Foo struct {
		Id    int
		Stuff string
	}
	DB.DropTable(&Foo{})
	if ok := DB.HasTable(&Foo{}); ok {
		t.Errorf("Table should not exist, but does")
	}
	if err := DB.CreateTable(&Foo{}).Error; err != nil {
		t.Errorf("Table should be created")
	}
	if ok := DB.HasTable(&Foo{}); !ok {
		t.Errorf("Table should exist, but HasTable informs it does not")
	}
}

func TestTableName(t *testing.T) {
	DB := DB.Model("")
	if DB.NewScope(Order{}).TableName() != "orders" {
		t.Errorf("Order's table name should be orders")
	}

	if DB.NewScope(&Order{}).TableName() != "orders" {
		t.Errorf("&Order's table name should be orders")
	}

	if DB.NewScope([]Order{}).TableName() != "orders" {
		t.Errorf("[]Order's table name should be orders")
	}

	if DB.NewScope(&[]Order{}).TableName() != "orders" {
		t.Errorf("&[]Order's table name should be orders")
	}

	DB.SingularTable(true)
	if DB.NewScope(Order{}).TableName() != "order" {
		t.Errorf("Order's singular table name should be order")
	}

	if DB.NewScope(&Order{}).TableName() != "order" {
		t.Errorf("&Order's singular table name should be order")
	}

	if DB.NewScope([]Order{}).TableName() != "order" {
		t.Errorf("[]Order's singular table name should be order")
	}

	if DB.NewScope(&[]Order{}).TableName() != "order" {
		t.Errorf("&[]Order's singular table name should be order")
	}

	if DB.NewScope(&Cart{}).TableName() != "shopping_cart" {
		t.Errorf("&Cart's singular table name should be shopping_cart")
	}

	if DB.NewScope(Cart{}).TableName() != "shopping_cart" {
		t.Errorf("Cart's singular table name should be shopping_cart")
	}

	if DB.NewScope(&[]Cart{}).TableName() != "shopping_cart" {
		t.Errorf("&[]Cart's singular table name should be shopping_cart")
	}

	if DB.NewScope([]Cart{}).TableName() != "shopping_cart" {
		t.Errorf("[]Cart's singular table name should be shopping_cart")
	}
	DB.SingularTable(false)
}

func TestSqlNullValue(t *testing.T) {
	DB.DropTable(&NullValue{})
	DB.AutoMigrate(&NullValue{})

	if err := DB.Save(&NullValue{Name: sql.NullString{String: "hello", Valid: true},
		Age:     sql.NullInt64{Int64: 18, Valid: true},
		Male:    sql.NullBool{Bool: true, Valid: true},
		Height:  sql.NullFloat64{Float64: 100.11, Valid: true},
		AddedAt: NullTime{Time: time.Now(), Valid: true},
	}).Error; err != nil {
		t.Errorf("Not error should raise when test null value")
	}

	var nv NullValue
	DB.First(&nv, "name = ?", "hello")

	if nv.Name.String != "hello" || nv.Age.Int64 != 18 || nv.Male.Bool != true || nv.Height.Float64 != 100.11 || nv.AddedAt.Valid != true {
		t.Errorf("Should be able to fetch null value")
	}

	if err := DB.Save(&NullValue{Name: sql.NullString{String: "hello-2", Valid: true},
		Age:     sql.NullInt64{Int64: 18, Valid: false},
		Male:    sql.NullBool{Bool: true, Valid: true},
		Height:  sql.NullFloat64{Float64: 100.11, Valid: true},
		AddedAt: NullTime{Time: time.Now(), Valid: false},
	}).Error; err != nil {
		t.Errorf("Not error should raise when test null value")
	}

	var nv2 NullValue
	DB.First(&nv2, "name = ?", "hello-2")
	if nv2.Name.String != "hello-2" || nv2.Age.Int64 != 0 || nv2.Male.Bool != true || nv2.Height.Float64 != 100.11 || nv2.AddedAt.Valid != false {
		t.Errorf("Should be able to fetch null value")
	}

	if err := DB.Save(&NullValue{Name: sql.NullString{String: "hello-3", Valid: false},
		Age:     sql.NullInt64{Int64: 18, Valid: false},
		Male:    sql.NullBool{Bool: true, Valid: true},
		Height:  sql.NullFloat64{Float64: 100.11, Valid: true},
		AddedAt: NullTime{Time: time.Now(), Valid: false},
	}).Error; err == nil {
		t.Errorf("Can't save because of name can't be null")
	}
}

func TestTransaction(t *testing.T) {
	tx := DB.Begin()
	u := User{Name: "transcation"}
	if err := tx.Save(&u).Error; err != nil {
		t.Errorf("No error should raise")
	}

	if err := tx.First(&User{}, "name = ?", "transcation").Error; err != nil {
		t.Errorf("Should find saved record")
	}

	if sqlTx, ok := tx.CommonDB().(*sql.Tx); !ok || sqlTx == nil {
		t.Errorf("Should return the underlying sql.Tx")
	}

	tx.Rollback()

	if err := tx.First(&User{}, "name = ?", "transcation").Error; err == nil {
		t.Errorf("Should not find record after rollback")
	}

	tx2 := DB.Begin()
	u2 := User{Name: "transcation-2"}
	if err := tx2.Save(&u2).Error; err != nil {
		t.Errorf("No error should raise")
	}

	if err := tx2.First(&User{}, "name = ?", "transcation-2").Error; err != nil {
		t.Errorf("Should find saved record")
	}

	tx2.Commit()

	if err := DB.First(&User{}, "name = ?", "transcation-2").Error; err != nil {
		t.Errorf("Should be able to find committed record")
	}
}

func TestRow(t *testing.T) {
	user1 := User{Name: "RowUser1", Age: 1, Birthday: now.MustParse("2000-1-1")}
	user2 := User{Name: "RowUser2", Age: 10, Birthday: now.MustParse("2010-1-1")}
	user3 := User{Name: "RowUser3", Age: 20, Birthday: now.MustParse("2020-1-1")}
	DB.Save(&user1).Save(&user2).Save(&user3)

	row := DB.Table("users").Where("name = ?", user2.Name).Select("age").Row()
	var age int64
	row.Scan(&age)
	if age != 10 {
		t.Errorf("Scan with Row")
	}
}

func TestRows(t *testing.T) {
	user1 := User{Name: "RowsUser1", Age: 1, Birthday: now.MustParse("2000-1-1")}
	user2 := User{Name: "RowsUser2", Age: 10, Birthday: now.MustParse("2010-1-1")}
	user3 := User{Name: "RowsUser3", Age: 20, Birthday: now.MustParse("2020-1-1")}
	DB.Save(&user1).Save(&user2).Save(&user3)

	rows, err := DB.Table("users").Where("name = ? or name = ?", user2.Name, user3.Name).Select("name, age").Rows()
	if err != nil {
		t.Errorf("Not error should happen, but got")
	}

	count := 0
	for rows.Next() {
		var name string
		var age int64
		rows.Scan(&name, &age)
		count++
	}
	if count != 2 {
		t.Errorf("Should found two records with name 3")
	}
}

func TestScan(t *testing.T) {
	user1 := User{Name: "ScanUser1", Age: 1, Birthday: now.MustParse("2000-1-1")}
	user2 := User{Name: "ScanUser2", Age: 10, Birthday: now.MustParse("2010-1-1")}
	user3 := User{Name: "ScanUser3", Age: 20, Birthday: now.MustParse("2020-1-1")}
	DB.Save(&user1).Save(&user2).Save(&user3)

	type result struct {
		Name string
		Age  int
	}

	var res result
	DB.Table("users").Select("name, age").Where("name = ?", user3.Name).Scan(&res)
	if res.Name != user3.Name {
		t.Errorf("Scan into struct should work")
	}

	var doubleAgeRes result
	DB.Table("users").Select("age + age as age").Where("name = ?", user3.Name).Scan(&doubleAgeRes)
	if doubleAgeRes.Age != res.Age*2 {
		t.Errorf("Scan double age as age")
	}

	var ress []result
	DB.Table("users").Select("name, age").Where("name in (?)", []string{user2.Name, user3.Name}).Scan(&ress)
	if len(ress) != 2 || ress[0].Name != user2.Name || ress[1].Name != user3.Name {
		t.Errorf("Scan into struct map")
	}
}

func TestRaw(t *testing.T) {
	user1 := User{Name: "ExecRawSqlUser1", Age: 1, Birthday: now.MustParse("2000-1-1")}
	user2 := User{Name: "ExecRawSqlUser2", Age: 10, Birthday: now.MustParse("2010-1-1")}
	user3 := User{Name: "ExecRawSqlUser3", Age: 20, Birthday: now.MustParse("2020-1-1")}
	DB.Save(&user1).Save(&user2).Save(&user3)

	type result struct {
		Name  string
		Email string
	}

	var ress []result
	DB.Raw("SELECT name, age FROM users WHERE name = ? or name = ?", user2.Name, user3.Name).Scan(&ress)
	if len(ress) != 2 || ress[0].Name != user2.Name || ress[1].Name != user3.Name {
		t.Errorf("Raw with scan")
	}

	rows, _ := DB.Raw("select name, age from users where name = ?", user3.Name).Rows()
	count := 0
	for rows.Next() {
		count++
	}
	if count != 1 {
		t.Errorf("Raw with Rows should find one record with name 3")
	}

	DB.Exec("update users set name=? where name in (?)", "jinzhu", []string{user1.Name, user2.Name, user3.Name})
	if DB.Where("name in (?)", []string{user1.Name, user2.Name, user3.Name}).First(&User{}).Error != gorm.RecordNotFound {
		t.Error("Raw sql to update records")
	}
}

func TestGroup(t *testing.T) {
	rows, err := DB.Select("name").Table("users").Group("name").Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			rows.Scan(&name)
		}
	} else {
		t.Errorf("Should not raise any error")
	}
}

func TestJoins(t *testing.T) {
	type result struct {
		Name  string
		Email string
	}

	user := User{
		Name:   "joins",
		Emails: []Email{{Email: "join1@example.com"}, {Email: "join2@example.com"}},
	}
	DB.Save(&user)

	var results []result
	DB.Table("users").Select("name, email").Joins("left join emails on emails.user_id = users.id").Where("name = ?", "joins").Scan(&results)
	if len(results) != 2 || results[0].Email != "join1@example.com" || results[1].Email != "join2@example.com" {
		t.Errorf("Should find all two emails with Join")
	}
}

func TestHaving(t *testing.T) {
	rows, err := DB.Select("name, count(*) as total").Table("users").Group("name").Having("name IN (?)", []string{"2", "3"}).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			var total int64
			rows.Scan(&name, &total)

			if name == "2" && total != 1 {
				t.Errorf("Should have one user having name 2")
			}
			if name == "3" && total != 2 {
				t.Errorf("Should have two users having name 3")
			}
		}
	} else {
		t.Errorf("Should not raise any error")
	}
}

func DialectHasTzSupport() bool {
	// NB: mssql and FoundationDB do not support time zones.
	if dialect := os.Getenv("GORM_DIALECT"); dialect == "mssql" || dialect == "foundation" {
		return false
	}
	return true
}

func TestTimeWithZone(t *testing.T) {
	var format = "2006-01-02 15:04:05 -0700"
	var times []time.Time
	GMT8, _ := time.LoadLocation("Asia/Shanghai")
	times = append(times, time.Date(2013, 02, 19, 1, 51, 49, 123456789, GMT8))
	times = append(times, time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.UTC))

	for index, vtime := range times {
		name := "time_with_zone_" + strconv.Itoa(index)
		user := User{Name: name, Birthday: vtime}

		if !DialectHasTzSupport() {
			// If our driver dialect doesn't support TZ's, just use UTC for everything here.
			user.Birthday = vtime.UTC()
		}

		DB.Save(&user)
		expectedBirthday := "2013-02-18 17:51:49 +0000"
		foundBirthday := user.Birthday.UTC().Format(format)
		if foundBirthday != expectedBirthday {
			t.Errorf("User's birthday should not be changed after save for name=%s, expected bday=%+v but actual value=%+v", name, expectedBirthday, foundBirthday)
		}

		var findUser, findUser2, findUser3 User
		DB.First(&findUser, "name = ?", name)
		foundBirthday = findUser.Birthday.UTC().Format(format)
		if foundBirthday != expectedBirthday {
			t.Errorf("User's birthday should not be changed after find for name=%s, expected bday=%+v but actual value=%+v or %+v", name, expectedBirthday, foundBirthday)
		}

		if DB.Where("id = ? AND birthday >= ?", findUser.Id, user.Birthday.Add(-time.Minute)).First(&findUser2).RecordNotFound() {
			t.Errorf("User should be found")
		}

		if !DB.Where("id = ? AND birthday >= ?", findUser.Id, user.Birthday.Add(time.Minute)).First(&findUser3).RecordNotFound() {
			t.Errorf("User should not be found")
		}
	}
}

func TestHstore(t *testing.T) {
	type Details struct {
		Id   int64
		Bulk gorm.Hstore
	}

	if dialect := os.Getenv("GORM_DIALECT"); dialect != "postgres" {
		t.Skip()
	}

	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS hstore").Error; err != nil {
		fmt.Println("\033[31mHINT: Must be superuser to create hstore extension (ALTER USER gorm WITH SUPERUSER;)\033[0m")
		panic(fmt.Sprintf("No error should happen when create hstore extension, but got %+v", err))
	}

	DB.Exec("drop table details")

	if err := DB.CreateTable(&Details{}).Error; err != nil {
		panic(fmt.Sprintf("No error should happen when create table, but got %+v", err))
	}

	bankAccountId, phoneNumber, opinion := "123456", "14151321232", "sharkbait"
	bulk := map[string]*string{
		"bankAccountId": &bankAccountId,
		"phoneNumber":   &phoneNumber,
		"opinion":       &opinion,
	}
	d := Details{Bulk: bulk}
	DB.Save(&d)

	var d2 Details
	if err := DB.First(&d2).Error; err != nil {
		t.Errorf("Got error when tried to fetch details: %+v", err)
	}

	for k := range bulk {
		if r, ok := d2.Bulk[k]; ok {
			if res, _ := bulk[k]; *res != *r {
				t.Errorf("Details should be equal")
			}
		} else {
			t.Errorf("Details should be existed")
		}
	}
}

func TestSetAndGet(t *testing.T) {
	if value, ok := DB.Set("hello", "world").Get("hello"); !ok {
		t.Errorf("Should be able to get setting after set")
	} else {
		if value.(string) != "world" {
			t.Errorf("Setted value should not be changed")
		}
	}

	if _, ok := DB.Get("non_existing"); ok {
		t.Errorf("Get non existing key should return error")
	}
}

func TestCompatibilityMode(t *testing.T) {
	DB, _ := gorm.Open("testdb", "")
	testdb.SetQueryFunc(func(query string) (driver.Rows, error) {
		columns := []string{"id", "name", "age"}
		result := `
		1,Tim,20
		2,Joe,25
		3,Bob,30
		`
		return testdb.RowsFromCSVString(columns, result), nil
	})

	var users []User
	DB.Find(&users)
	if (users[0].Name != "Tim") || len(users) != 3 {
		t.Errorf("Unexcepted result returned")
	}
}

func TestOpenExistingDB(t *testing.T) {
	DB.Save(&User{Name: "jnfeinstein"})
	dialect := os.Getenv("GORM_DIALECT")

	db, err := gorm.Open(dialect, DB.DB())
	if err != nil {
		t.Errorf("Should have wrapped the existing DB connection")
	}

	var user User
	if db.Where("name = ?", "jnfeinstein").First(&user).Error == gorm.RecordNotFound {
		t.Errorf("Should have found existing record")
	}
}

func BenchmarkGorm(b *testing.B) {
	b.N = 2000
	for x := 0; x < b.N; x++ {
		e := strconv.Itoa(x) + "benchmark@example.org"
		email := BigEmail{Email: e, UserAgent: "pc", RegisteredAt: time.Now()}
		// Insert
		DB.Save(&email)
		// Query
		DB.First(&BigEmail{}, "email = ?", e)
		// Update
		DB.Model(&email).UpdateColumn("email", "new-"+e)
		// Delete
		DB.Delete(&email)
	}
}

func BenchmarkRawSql(b *testing.B) {
	DB, _ := sql.Open("postgres", "user=gorm DB.ame=gorm sslmode=disable")
	DB.SetMaxIdleConns(10)
	insertSql := "INSERT INTO emails (user_id,email,user_agent,registered_at,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"
	querySql := "SELECT * FROM emails WHERE email = $1 ORDER BY id LIMIT 1"
	updateSql := "UPDATE emails SET email = $1, updated_at = $2 WHERE id = $3"
	deleteSql := "DELETE FROM orders WHERE id = $1"

	b.N = 2000
	for x := 0; x < b.N; x++ {
		var id int64
		e := strconv.Itoa(x) + "benchmark@example.org"
		email := BigEmail{Email: e, UserAgent: "pc", RegisteredAt: time.Now()}
		// Insert
		DB.QueryRow(insertSql, email.UserId, email.Email, email.UserAgent, email.RegisteredAt, time.Now(), time.Now()).Scan(&id)
		// Query
		rows, _ := DB.Query(querySql, email.Email)
		rows.Close()
		// Update
		DB.Exec(updateSql, "new-"+e, time.Now(), id)
		// Delete
		DB.Exec(deleteSql, id)
	}
}
