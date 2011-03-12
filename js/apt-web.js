var urls = [];
function cmirror(obj) {
	if (disabled) { return; }
	var i;

	if (urls.length == 0) {
		var o = $('#urls li a');
		for (i=0; i<o.length; i++) {
			urls[i] = $(o[i]).attr('href').replace(base_url, '');
		}
	}

	var new_url = mirrors[$(obj).attr('value')];
	var txt = '';
	var url;

	for (i=0; i<urls.length; i++) {
		url = new_url + urls[i];
		txt += '<li><a href="'+url+'">'+url+'</a></li>';
	}

	$('#urls ul').html(txt);
}

$(document).ready(function() {
    // Add Colorbox
    $('#pkg-newest, #pkg-extra, #pkg-rec, #pkg-suggest, #pkg-inst')
        .find('a')
        .colorbox({'maxWidth': 800, 'maxHeight': '75%'});

    // Add ZeroClipboard
    var urls_head = $('#urls h2');
    if (urls_head.length > 0) {
        var copy = $('<span id="copy-to-clipboard">copy to clipboard</span>');
        urls_head.wrap($('<div id="urls-head"/>'));
        urls_head.after(copy);
        urls_head.after(' &mdash; ');

        var clip = new window.ZeroClipboard.Client();
        clip.setText('');
        clip.setHandCursor(true);
        clip.setCSSEffects(true);
        clip.addEventListener('complete', function(client, text) {
            copy.html('copied!');
            setTimeout(function() { copy.html('copy to clipboard') }, 2000);
        });
        clip.addEventListener('mouseDown', function(client) {
            var urls = [];
            $('#urls ul li a').each(function() {
                urls.push($(this).attr('href'));
            });
            urls = urls.join("\n");
            clip.setText(urls);
        });
        clip.glue('copy-to-clipboard');
    }
});

