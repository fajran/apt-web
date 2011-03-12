<?

include('config.php');
include('common.php');

$dists = $_repo_list;
$mirrors = $_mirror_list;

$packages = '';
$dist = 0;
$mirror = 0;

if (isset($_POST['submit'])) {

	$packages = isset($_POST['packages']) ? trim($_POST['packages']) : '';
	$dist = isset($_POST['dist']) ? intval(trim($_POST['dist'])) : 0;
	$mirror = isset($_POST['mirror']) ? intval(trim($_POST['mirror'])) : 0;

	if (($packages != '') && (isset($_repo_list[$dist]))) {

		// Get package dependencies and their URLs
		$res = apt_install($_repo_list[$dist][0], $packages);
		$res = parse_install($res);
		$extra = &$res['extra'];
		$suggested = &$res['suggested'];
		$recommended = &$res['recommended'];
		$install = &$res['install'];
		$tbInstalled = &$res['packages'];
		$newest = &$res['newest'];

		if (isset($_mirror_list[$mirror])) {
			while (list($key, $val) = each($tbInstalled)) {
				$tbInstalled[$key][0] = convert_url($tbInstalled[$key][0], $_mirror_list[$mirror][0]);
			}
		}
	}

}

function description_link($pkg) {
	global $dist;
	return "<a href=\"desc.php?p=$pkg&d=$dist&width=650&height=550\" class=\"thickbox\" title=\"$pkg\">$pkg</a>";
}

?>
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0//EN" "http://www.w3.org/TR/REC-html40/strict.dtd">
<html><head>
<title>Which files should I download?</title>
<style type="text/css">@import "css/thickbox.css";</style>
<style type="text/css">@import "css/style.css";</style>
<script type="text/javascript" src="js/jquery.min.js"></script>
<script type="text/javascript" src="js/thickbox.js"></script>
<script type="text/javascript">
var disabled = true;
<? if (isset($_POST['submit'])) { ?>
disabled = false;
var base_url = "<?=rtrim($_mirror_list[$mirror][0], '/');?>";
<? } ?>
var mirrors = [];
<?
reset($_mirror_list);
while (list($key, $val) = each($_mirror_list)) {
	print('mirrors['.$key.'] = "'.rtrim($val[0], '/').'";'."\n");
}
?>
</script>
<script type="text/javascript" src="js/apt-web.js"></script>
</head><body>

<h1>Which files should I download?</h1>

<div id="form">
<form method="post">

<p><label>Base distribution</label>
	<select name="dist" value="<?=$dist;?>">
<? while (list($key, $val) = each($dists)) { ?>
		<option <?=($dist==$key?'selected="selected" ':' ');?>value="<?=$key;?>"><?=$val[1];?></option>
<? } ?>
	</select>
</p>

<p><label>Mirror</label>
	<select name="mirror" onchange="cmirror(this)">
<? while (list($key, $val) = each($mirrors)) { ?>
		<option <?=($mirror==$key?'selected="selected" ':' ');?>value="<?=$key;?>"><?=$val[1];?></option>
<? } ?>
	</select>
</p>

<p><label>Packages</label>
	<input class="txt" type="text" name="packages" value="<?=$packages;?>"/>
</p>

<p><input type="submit" name="submit" value="submit"/></p>

</form>
</div>

<? if (isset($_POST['submit'])) { ?>

<div id="left">

<? if (!empty($newest)) { ?>
<h2>Already Installed</h2>
<p id="pkg-newest"><?=implode(', ', array_map('description_link', $newest));?></p>
<? } ?>

<? if (!empty($extra)) { ?>
<h2>Extra</h2>
<p id="pkg-extra"><?=implode(', ', array_map('description_link', $extra));?></p>
<? } ?>

<? if (!empty($recommended)) { ?>
<h2>Recommended</h2>
<p id="pkg-rec"><?=implode(', ', array_map('description_link', $recommended));?></p>
<? } ?>

<? if (!empty($suggested)) { ?>
<h2>Suggested</h2>
<p id="pkg-suggest"><?=implode(', ', array_map('description_link', $suggested));?></p>
<? } ?>

<? if (!empty($install)) { ?>
<h2>To Be Installed</h2>
<p id="pkg-inst"><?=implode(', ', array_map('description_link', $install));?></p>

</div>
<div id="right">

<div id="urls">
<h2>URLs</h2>
<ul>
<? foreach ($tbInstalled as $package) { ?>
<li><a href="<?=$package[0];?>"><?=$package[0];?></a></li>
<? } ?>
</ul>
</div>

</div>
<? } ?>

</div>

<? } ?>

<div id="footer"><a href="http://github.com/fajran/apt-web/">apt-web</a> - copyright &copy; 2007 &ndash; 2011 - <a href="http://fajran.web.id">Fajran Iman Rusadi</a></div>

</body></html>
