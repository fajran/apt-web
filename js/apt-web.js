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

