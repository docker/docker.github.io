// Right nav highlighting
var sidebarObj = (document.getElementsByClassName("sidebar")[0]) ? document.getElementsByClassName("sidebar")[0] : document.getElementsByClassName("sidebar-home")[0];

// ensure that the left nav visibly displays the current topic
var current = document.getElementsByClassName("active currentPage");
var body = document.getElementsByClassName("col-content content");
if (current[0]) {
    if (sidebarObj) {
        current[0].scrollIntoView(true);
        body[0].scrollIntoView(true);
    }
    // library hack
    if (document.location.pathname.indexOf("/samples/") > -1) {
        $(".currentPage").closest("ul").addClass("in");
    }
}

function navClicked(sourceLink) {
    var classString = document.getElementById("#item" + sourceLink).className;
    if (classString.indexOf(" in") > -1) {
        //collapse
        document.getElementById("#item" + sourceLink).className = classString.replace(" in", "");
    } else {
        //expand
        document.getElementById("#item" + sourceLink).className = classString.concat(" in");
    }
}

var outputLetNav = [];
var totalTopics = 0;

function pageIsInSection(tree) {
    function processBranch(branch) {
        for (let k = 0; k < branch.length; k++) {
            if (branch[k].section) {
                processBranch(branch[k].section);
            } else {
                if (branch[k].path === pageURL && !branch[k].nosync) {
                    found = true;
                    break;
                }
            }
        }
    }

    let found = false;
    processBranch(tree);
    return found;
}

function walkTree(tree) {
    for (let j = 0; j < tree.length; j++) {
        totalTopics++;
        if (tree[j].section) {
            let sectionHasPath = pageIsInSection(tree[j].section);
            outputLetNav.push('<li><a onclick="navClicked(' + totalTopics + ')" data-target="#item' + totalTopics + '" data-toggle="collapse" data-parent="#stacked-menu"')
            if (sectionHasPath) {
                outputLetNav.push('aria-expanded="true"')
            } else {
                outputLetNav.push('class="collapsed" aria-expanded="false"')
            }
            outputLetNav.push(">" + tree[j].sectiontitle + '<span class="caret arrow"></span></a>');
            outputLetNav.push('<ul class="nav collapse');
            if (sectionHasPath) outputLetNav.push(" in");
            outputLetNav.push('" id="#item' + totalTopics + '" aria-expanded="');
            if (sectionHasPath) {
                outputLetNav.push("true");
            } else {
                outputLetNav.push("false");
            }
            outputLetNav.push('">');
            walkTree(tree[j].section);
            outputLetNav.push("</ul></li>");
        } else {
            // just a regular old topic; this is a leaf, not a branch; render a link!
            outputLetNav.push('<li><a href="' + tree[j].path + '"')
            if (tree[j].path === pageURL && !tree[j].nosync) {
                outputLetNav.push('class="active currentPage"')
            }
            outputLetNav.push(">" + tree[j].title + "</a></li>")
        }
    }
}

function renderNav(docstoc) {
    for (let i = 0; i < docstoc.horizontalnav.length; i++) {
        if (docstoc.horizontalnav[i].path === pageURL || pageIsInSection(docstoc[docstoc.horizontalnav[i].node])) {
            // This is the current section. Set the corresponding header-nav link
            // to active, and build the left-hand (vertical) navigation
            document.getElementById(docstoc.horizontalnav[i].node).closest("li").classList.add("active")
            walkTree(docstoc[docstoc.horizontalnav[i].node]);
            document.getElementById("jsTOCLeftNav").innerHTML = outputLetNav.join("");
        }
    }
    // Scroll the current menu item into view. We actually pick the item *above*
    // the current item to give some headroom above
    scrollMenuItem("#jsTOCLeftNav a.currentPage")
}

// Scroll the given menu item into view. We actually pick the item *above*
// the current item to give some headroom above
function scrollMenuItem(selector) {
    let item = document.querySelector(selector)
    if (item) {
        item = item.closest("li")
    }
    if (item) {
        item = item.previousElementSibling
    }
    if (item) {
        item.scrollIntoView(true)
        if (window.location.hash.length < 2) {
            // Scrolling the side-navigation may scroll the whole page as well
            // this is a dirty hack to scroll the main content back to the top
            // if we're not on a specific anchor
            document.querySelector("main.col-content").scrollIntoView(true)
        }
    }
}

function highlightRightNav(heading) {
    $("#my_toc a.active").removeClass("active");

    if (heading !== "title") {
        $("#my_toc a[href='#" + heading + "']").addClass("active");
    }
}

var currentHeading = "";
$(window).scroll(function () {
    var headingPositions = [];
    $("h1, h2, h3, h4, h5, h6").each(function () {
        if (this.id === "") this.id = "title";
        headingPositions[this.id] = this.getBoundingClientRect().top;
    });
    headingPositions.sort();
    // the headings have all been grabbed and sorted in order of their scroll
    // position (from the top of the page). First one is toppermost.
    for (var key in headingPositions) {
        if (!headingPositions.hasOwnProperty(key)) {
            continue;
        }
        if (headingPositions[key] > 0 && headingPositions[key] < 200) {
            if (currentHeading !== key) {
                // a new heading has scrolled to within 200px of the top of the page.
                // highlight the right-nav entry and de-highlight the others.
                highlightRightNav(key);
                currentHeading = key;
            }
            break;
        }
    }
});

/*
 *
 * toggle menu *********************************************************************
 *
 */

$("#menu-toggle").click(function (e) {
    e.preventDefault();
    $(".wrapper").toggleClass("right-open");
    $(".col-toc").toggleClass("col-toc-hidden");
});

$("#menu-toggle-left").click(function (e) {
    e.preventDefault();
    $(".col-nav").toggleClass("col-toc-hidden");
});

$(".navbar-toggle").click(function () {
    $("#sidebar-nav").each(function () {
        $(this).toggleClass("hidden-sm");
        $(this).toggleClass("hidden-xs");
    });
});

var navHeight = $(".navbar").outerHeight(true) + 80;

$(document.body).scrollspy({
    target: "#leftCol",
    offset: navHeight
});

function loadHash(hashObj) {
    // Using jQuery's animate() method to add smooth page scroll
    // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
    $("html, body").animate({scrollTop: $(hashObj).offset().top - 80}, 800);
}

$(document).ready(function () {
    // Add smooth scrolling to all links
    $(".toc-nav a").on("click", function (event) {
        // Make sure this.hash has a value before overriding default behavior
        if (this.hash !== "") {
            // Prevent default anchor click behavior
            event.preventDefault();

            // Store hash
            var hash = this.hash;
            loadHash(hash);

            // Add hash (#) to URL when done scrolling (default click behavior)
            window.location.hash = hash;
        }
    });
    if (window.location.hash) loadHash(window.location.hash);
});


$(document).ready(function () {
    $(".sidebar").Stickyfill();

    // Add smooth scrolling to all links
    $(".nav-sidebar ul li a").on("click", function (event) {

        // Make sure this.hash has a value before overriding default behavior
        if (this.hash !== "") {
            // Prevent default anchor click behavior
            event.preventDefault();

            // Store hash
            var hash = this.hash;

            // Using jQuery's animate() method to add smooth page scroll
            // The optional number (800) specifies the number of milliseconds it takes to scroll to the specified area
            $("html, body").animate({
                scrollTop: $(hash).offset().top - 80
            }, 800, function () {
                // Add hash (#) to URL when done scrolling (default click behavior)
                window.location.hash = hash;
            });
        }
    });
});


/*
 *
 * make dropdown show on hover *********************************************************************
 *
 */

$("ul.nav li.dropdown").hover(function () {
    $(this).find(".dropdown-menu").stop(true, true).delay(200).fadeIn(500);
}, function () {
    $(this).find(".dropdown-menu").stop(true, true).delay(200).fadeOut(500);
});


/*
 *
 * TEMP HACK For side menu*********************************************************************
 *
 */

$(".nav-sidebar ul li a").click(function () {
    $(this).addClass("collapse").siblings().toggleClass("in");
});

if ($(".nav-sidebar ul a.active").length !== 0) {
    $(".nav-sidebar ul").click(function () {
        $(this).addClass("collapse in").siblings;
    });
}

/*
 *
 * Components *********************************************************************
 *
 */

$(function () {
    $('[data-toggle="tooltip"]').tooltip()
});

// sync tabs with the same data-group
window.onload = function () {
    $(".nav-tabs > li > a").click(function (e) {
        var group = $(this).attr("data-group");
        $('.nav-tabs > li > a[data-group="' + group + '"]').tab("show");
    });
};
