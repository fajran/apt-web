<?

include('config.php');
include('common.php');

$package = isset($_GET['p']) ? parse($_GET['p']) : '';
$dist = isset($_GET['d']) ? trim($_GET['d']) : '';

$valid = false;

if (($package != '') && (isset($_repo_list[$dist]))) {

	// Get package description
	$data = apt_show($_repo_list[$dist][0], $package);

	$line = true;
	$info = array();
	while ($line) {

		// Long description
		if (strpos($line, 'Description') === 0) {
			$short_desc = substr($line, 13);
			$long_desc = array();

			$line = array_shift($data);
			while (strpos($line, ' ') === 0) {
				if ($line == ' .') {
                    $long_desc[] = '';
				}
				else {
					$long_desc[] = trim($line);
				}
				$line = array_shift($data);
			}

			$long_desc = implode("\n", $long_desc);
		}

		// Other informations
		else {
			if ($line !== true) {
				$info[] = explode(': ', $line);
			}

			$line = array_shift($data);
		}

	}
	
	$valid = true;
}

?>
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0//EN" "http://www.w3.org/TR/REC-html40/strict.dtd">
<html><head>
<title><?=$package;?></title>
<style type="text/css">@import "css/desc.css";</style>
</head><body>
<? if ($valid) { ?>

<div id="desc">
<h1><?=$package;?></h1>

<p id="short"><?=$short_desc;?></p>

<div id="long">
<?=$long_desc;?>
</div>

<div id="info">
<table>
<? foreach ($info as $d) { ?>
<tr><td><?=$d[0];?></td><td><?=htmlentities($d[1]);?></td></tr>
<? } ?>
</table>
</div>
</div>

<? } else { ?>

<h1>Invalid input</h1>

<? } ?>
</body></html>
