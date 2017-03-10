/* Obly run this if we are online*/
if (window.navigator.onLine) {
  var dockerVersion = '1.7';
  /* This JSON file contains a current list of all docs versions of Docker */
  $.getJSON("https://docs.docker.com/js/archives.json", function(result){
    var outerDivStart = '<div style="padding-top: 10px; padding-bottom: 10px; min-height: 34px; border: 1px solid #A3733F; background-color: #FFE1C0; color: #A3733F"><div class="container"><div style="text-align: center"><span id="archive-list">This is <b><a href="https://docs.docker.com/docsarchive/" style="color: #A3733F; text-decoration: underline !important">archived documentation</a></b> for Docker&nbsp;' + dockerVersion + '. Go to the <a style="color: #A3733F; text-decoration: underline !important" href="https://docs.docker.com/">latest docs</a> or a different version:&nbsp;&nbsp;</span>' +
                               '<span style="z-index: 1001" class="dropdown">';
    var listStart = '<ul class="dropdown-menu" role="menu" aria-labelledby="archive-menu">';
    var listEnd = '</ul>';
    var outerDivEnd = '</span></div></div></div>';
    var buttonCode = null;
    var listItems = new Array();
    $.each(result, function(i, field){
      var prettyName = 'Docker ' + field.name.replace("v", "");
      // If this archive has current = true, and we don't already have a button
      if ( field.current && buttonCode == null ) {
        // Get the button code
        buttonCode = '<button id="archive-menu" data-toggle="dropdown" class="btn dropdown-toggle" style="border: 1px solid #A3733F; background-color: #fff; color: #A3733F;">' + prettyName + '&nbsp;(current) &nbsp;<span class="caret"></span></button>';
        // The link is different for the current release
        listItems.push('<li role="presentation"><a role="menuitem" tabindex="-1" href="https://docs.docker.com/">' + prettyName + '</a></li>');
      } else {
        listItems.push('<li role="presentation"><a role="menuitem" tabindex="-1" href="https://docs.docker.com/' + field.name + '/">' + prettyName + '</a></li>');
      }
    });
    $( 'body' ).prepend(outerDivStart + buttonCode + listStart + listItems.join("") + listEnd + outerDivEnd);
  });
}
