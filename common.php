<?

function parse_clean($package) {
	$pattern = '/[^a-z0-9.-]/';
	$replacement = '';
	return preg_replace($pattern, $replacement, $package);
}

function parse($packages) {
	$p = explode(' ', $packages);
	$list = array();
	foreach ($p as $package) {
		$list[] = parse_clean($package);
	}
	return implode(' ', $list);
}

function get_dir($repo) {
	global $_repo_dir;
	
	$dir = $_repo_dir . $repo;

	if (is_dir($dir)) {
		return $dir;
	}
	else {
		return false;
	}
}

function run_apt($cmd, $repo) {	
	$dir = get_dir($repo);
	if ($dir !== false) {
		chdir($dir);
		exec($cmd, $output);
		return $output;
	}
	else {
		return false;
	}
}

function apt_search($repo, $packages) {
	global $_file_apt_cache;

	$cmd = $_file_apt_cache.' -c=apt.conf search '.parse($packages);
	return run_apt($cmd, $repo);
}

function apt_show($repo, $packages) {
	global $_file_apt_cache;

	$cmd = $_file_apt_cache.' -c=apt.conf show '.parse($packages);
	return run_apt($cmd, $repo);
}

function apt_install($repo, $packages) {
	global $_file_apt_get;

	$cmd = $_file_apt_get.' -c=apt.conf -y --print-uris install '.parse($packages);
	return run_apt($cmd, $repo);
}

function parse_install($data) {
	$extra = array();
	$suggested = array();
	$recommended = array();
	$install = array();
	$packages = array();
	$newest = array();

	$line = true;
	while ($line) {

		if (substr($line, -30) == 'is already the newest version.') {
			$newest[] = substr($line, 0, -31);
			$line = array_shift($data);
		}

		else if (strpos($line, 'The following extra packages will be installed:') === 0) {
			$line = array_shift($data);
			while (strpos($line, ' ') === 0) {
				$extra = array_merge($extra, explode(' ', trim($line)));
				$line = array_shift($data);
			}
		}

		else if (strpos($line, 'Suggested packages:') === 0) {
			$line = array_shift($data);
			while (strpos($line, ' ') === 0) {
				$suggested = array_merge($suggested, explode(' ', trim($line)));
				$line = array_shift($data);
			}
		}

		else if (strpos($line, 'Recommended packages:') === 0) {
			$line = array_shift($data);
			while (strpos($line, ' ') === 0) {
				$recommended = array_merge($recommended, explode(' ', trim($line)));
				$line = array_shift($data);
			}
		}

		else if (strpos($line, 'The following NEW packages will be installed:') === 0) {
			$line = array_shift($data);
			while (strpos($line, ' ') === 0) {
				$install = array_merge($install, explode(' ', trim($line)));
				$line = array_shift($data);
			}
		}

		else if (strpos($line, 'After unpacking') === 0) {
			$line = array_shift($data);
			while (strpos($line, '\'') === 0) {
				$info = explode(' ', $line);
				$info[0] = substr($info[0], 1, -1);
				$packages[] = $info;
				$line = array_shift($data);
			}
		}

		else {
			$line = array_shift($data);
		}
	
	}

	return array(
		'extra' => $extra,
		'suggested' => $suggested,
		'recommended' => $recommended,
		'install' => $install,
		'packages' => $packages,
		'newest' => $newest
	);
}

function convert_url($url, $mirror_url) {
	global $_repo_mirror_base;
	return str_replace($_repo_mirror_base, $mirror_url, $url);
}

?>
