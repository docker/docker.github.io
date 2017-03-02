// Right nav highlighting
var sidebarObj = (document.getElementsByClassName("sidebar")[0]) ? document.getElementsByClassName("sidebar")[0] : document.getElementsByClassName("sidebar-home")[0]
var sidebarBottom = sidebarObj.getBoundingClientRect().bottom;
var footerTop = document.getElementsByClassName("footer")[0].getBoundingClientRect().top;
var headerOffset = document.getElementsByClassName("container-fluid")[0].getBoundingClientRect().bottom;

// ensure that the left nav visibly displays the current topic
var current = document.getElementsByClassName("active currentPage");
if (current[0]) {
    if (sidebarObj) {
      current[0].scrollIntoView(true);
      sidebarObj.scrollTop -= 150;
    }
  }

function highlightRightNav(heading)
{
  if (heading == "title")
  {
    history.replaceState({},"Top of page on " + document.location.pathname,document.location.protocol +"//"+ document.location.hostname + (location.port ? ':'+location.port: '') + document.location.pathname);
    $("#my_toc a").each(function(){
      $(this).removeClass("active");
    });
    $("#sidebar-wrapper").animate({
      scrollTop: 0
    },800);
  } else {
    var targetAnchorHREF = document.location.protocol +"//"+ document.location.hostname + (location.port ? ':'+location.port: '') + document.location.pathname + "#" + heading;
    // make sure we aren't filtering out that heading level
    var noFilterFound = false;
    $("#my_toc a").each(function(){
      if (this.href==targetAnchorHREF) {
        noFilterFound = true;
      }
    });
    // now, highlight that heading
    if (noFilterFound)
    {
      $("#my_toc a").each(function(){
        //console.log("right-nav",this.href);
        if (this.href==targetAnchorHREF)
        {
          history.replaceState({},this.innerText,targetAnchorHREF);
          $(this).addClass("active");
          var sidebarOffset = (sidebarBottom > 200) ? 200 : headerOffset - 20;
          $("#sidebar-wrapper").animate({
            scrollTop: $("#sidebar-wrapper").scrollTop() + $(this).position().top - sidebarOffset
          },100);
          //document.getElementById("sidebar-wrapper").scrollTop = this.getBoundingClientRect().top - 200;
        } else {
          $(this).removeClass("active");
        }
      });
    }
  }
}
function checkNavSizes()
{
  sidebarBottom = sidebarObj.getBoundingClientRect().bottom;
  footerTop = document.getElementsByClassName("footer")[0].getBoundingClientRect().top;
  headerOffset = document.getElementsByClassName("container-fluid")[0].getBoundingClientRect().bottom;
  if (footerTop < sidebarBottom || (sidebarBottom < footerTop && sidebarBottom < $(window).height()))
  {

    // the footer is overlapping the sidebar
    var sidebarHeight = (footerTop < $(window).height()) ? footerTop : $(window).height();
    var tocNavHeight = (footerTop < $(window).height()) ? footerTop : $(window).height();
    sidebarObj.style.height = sidebarHeight + "px";
    document.getElementsByClassName("toc-nav")[0].style.height = tocNavHeight + "px";
    $(sidebarObj).clearQueue().finish();
    setTimeout(highlightRightNav(currentHeading),1);
  }
}
$(window).resize(function() {
  checkNavSizes();
});
var currentHeading = "";
$(window).scroll(function(){
  checkNavSizes();
  var headingPositions = new Array();
  $("h1, h2, h3, h4, h5, h6").each(function(){
    if (this.id == "") this.id="title";
    headingPositions[this.id]=this.getBoundingClientRect().top;
  });
  headingPositions.sort();
  // the headings have all been grabbed and sorted in order of their scroll
  // position (from the top of the page). First one is toppermost.
  for(var key in headingPositions)
  {
    if (headingPositions[key] > 0 && headingPositions[key] < 200)
    {
      if (currentHeading != key)
      {
        // a new heading has scrolled to within 200px of the top of the page.
        // highlight the right-nav entry and de-highlight the others.
        highlightRightNav(key);
        currentHeading = key;
      }
      break;
    }
  }
});


// Cookie functions
function createCookie(name,value,days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days*24*60*60*1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + value + expires + "; path=/";
}

function readCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
}

function eraseCookie(name) {
    createCookie(name,"",-1);
}
/*
if (readCookie("night") == "true") {
  document.getElementById('pagestyle').setAttribute('href', '/css/style-alt.css');
  $('#switch-style').prop('checked', true);
} else {
  document.getElementById('pagestyle').setAttribute('href', '/css/style.css');
  $('#switch-style').prop('checked', false);
}
*/


/*
 *
 * toggle menu *********************************************************************
 *
 */

$("#menu-toggle").click(function(e) {
        e.preventDefault();
        $(".wrapper").toggleClass("right-open");
    });
$(".navbar-toggle").click(function(){
  $(".sidebar, .sidebar-home").each(function(){
    if($(this).hasClass("hidden-sm")) {
      $(this).removeClass("hidden-sm");
      $(this).css("overflow-y","scroll");
    } else {
      $(this).addClass("hidden-sm");
    }
    if($(this).hasClass("hidden-xs")) {
      $(this).removeClass("hidden-xs");
      $(this).css("overflow-y","scroll");
    } else {
      $(this).addClass("hidden-xs");
    }
  });
});

var navHeight = $('.navbar').outerHeight(true) + 80;

$(document.body).scrollspy({
	target: '#leftCol',
	offset: navHeight
});

function loadHash(hashObj)
{
  // Using jQuery's animate() method to add smooth page scroll
  // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
  $('html, body').animate({
    scrollTop: $(hashObj).offset().top-80
  }, 800);
}

$(document).ready(function(){
  // Add smooth scrolling to all links
  // $( ".toc-nav a" ).addClass( "active" );
  $(".toc-nav a").on('click', function(event) {
    // $(this).addClass('active');
    // Make sure this.hash has a value before overriding default behavior
    if (this.hash !== "") {
      // Prevent default anchor click behavior
      event.preventDefault();

      // Store hash
      var hash = this.hash;
      loadHash(hash);

      // Add hash (#) to URL when done scrolling (default click behavior)
      window.location.hash = hash;

    } // End if
  });
  if (window.location.hash) loadHash(window.location.hash);
});


$(document).ready(function(){
  // Add smooth scrolling to all links
  $(".nav-sidebar ul li a").on('click', function(event) {

    // Make sure this.hash has a value before overriding default behavior
    if (this.hash !== "") {
      // Prevent default anchor click behavior
      event.preventDefault();

      // Store hash
      var hash = this.hash;

      // Using jQuery's animate() method to add smooth page scroll
      // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
      $('html, body').animate({
        scrollTop: $(hash).offset().top-80
      }, 800, function(){

        // Add hash (#) to URL when done scrolling (default click behavior)
        window.location.hash = hash;
      });
    } // End if
  });
});


/*
 *
 * make dropdown show on hover *********************************************************************
 *
 */

$('ul.nav li.dropdown').hover(function() {
  $(this).find('.dropdown-menu').stop(true, true).delay(200).fadeIn(500);
}, function() {
  $(this).find('.dropdown-menu').stop(true, true).delay(200).fadeOut(500);
});

/*
 *
 * swapStyleSheet*********************************************************************
 *
 */

function swapStyleSheet(sheet) {
    document.getElementById('pagestyle').setAttribute('href', sheet);
}


$('#switch-style').change(function() {

    if ($(this).is(':checked')) {
        swapStyleSheet('/css/style-alt.css');
        createCookie("night",true,999)
    } else {
        swapStyleSheet('/css/style.css');
        eraseCookie("night")
    }
});

/*
 *
 * TEMP HACK For side menu*********************************************************************
 *
 */

$('.nav-sidebar ul li a').click(function() {
    $(this).addClass('collapse').siblings().toggleClass('in');
});

if($('.nav-sidebar ul a.active').length != 0)
{
  $('.nav-sidebar ul').click(function() {
      $(this).addClass('collapse in').siblings;
  });
}
