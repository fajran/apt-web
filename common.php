<?

// Clean the package name
function parse_clean($package) {
	$pattern = '/[^a-z0-9.-]/';
	$replacement = '';
	return preg_replace($pattern, $replacement, $package);
}

// Parse given packages string
function parse($packages) {
	$p = explode(' ', $packages);
	$list = array();
	foreach ($p as $package) {
		$list[] = parse_clean($package);
	}
	return implode(' ', $list);
}

// Get repository's directory
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

// Run apt-get or apt-cache safely (cross your fingers :-)
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

// Search packages (not used yet!)
function apt_search($repo, $packages) {
	global $_file_apt_cache;

	$cmd = $_file_apt_cache.' -c=apt.conf search '.parse($packages);
	return run_apt($cmd, $repo);
}

// Get package description
function apt_show($repo, $packages) {
	global $_file_apt_cache;

	$cmd = $_file_apt_cache.' -c=apt.conf show '.parse($packages);
	return run_apt($cmd, $repo);
}

// Get URLs of ready to be installed packages
function apt_install($repo, $packages) {
	global $_file_apt_get;

	$cmd = $_file_apt_get.' -c=apt.conf -y --print-uris install '.parse($packages);
	return run_apt($cmd, $repo);
}

// Retrieve extra, suggested, recommended, and to be installed packages
function parse_install($data) {
	$extra = array();
	$suggested = array();
	$recommended = array();
	$install = array();
	$packages = array();
	$newest = array();

	$line = true;
	while ($line) {

		// Already installed package, skip!
		if (substr($line, -30) == 'is already the newest version.') {
			$newest[] = substr($line, 0, -31);
			$line = array_shift($data);
		}

		// Extra packages
		else if (strpos($line, 'The following extra packages will be installed:') === 0) {
			$line = array_shift($data);
			while (strpos($line, ' ') === 0) {
				$extra = array_merge($extra, explode(' ', trim($line)));
				$line = array_shift($data);
			}
		}

		// Suggested packages
		else if (strpos($line, 'Suggested packages:') === 0) {
			$line = array_shift($data);
			while (strpos($line, ' ') === 0) {
				$suggested = array_merge($suggested, explode(' ', trim($line)));
				$line = array_shift($data);
			}
		}

		// Recommended packages
		else if (strpos($line, 'Recommended packages:') === 0) {
			$line = array_shift($data);
			while (strpos($line, ' ') === 0) {
				$recommended = array_merge($recommended, explode(' ', trim($line)));
				$line = array_shift($data);
			}
		}

		// Ready to be installed (downloaded) packages
		else if (strpos($line, 'The following NEW packages will be installed:') === 0) {
			$line = array_shift($data);
			while (strpos($line, ' ') === 0) {
				$install = array_merge($install, explode(' ', trim($line)));
				$line = array_shift($data);
			}
		}

		else if ((strpos($line, 'After unpacking') === 0) ||
			(strpos($line, 'After this operation') === 0)) {

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

// Get URL from a mirror
function convert_url($url, $mirror_url) {
	global $_repo_mirror_base;
	$base = rtrim($_repo_mirror_base, '/');
	$mirror_url = rtrim($mirror_url, '/');
	return str_replace($base, $mirror_url, $url);
}

?>
