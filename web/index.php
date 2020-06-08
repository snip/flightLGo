<html>
<body>
<a href="https://github.com/snip/flightLGo" class="github-corner" aria-label="View source on GitHub"><svg width="80" height="80" viewBox="0 0 250 250" style="fill:#151513; color:#fff; position: absolute; top: 0; border: 0; left: 0; transform: scale(-1, 1);" aria-hidden="true"><path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z"></path><path d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2" fill="currentColor" style="transform-origin: 130px 106px;" class="octo-arm"></path><path d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1,77.6 C213.8,80.2 216.3,84.9 216.3,84.9 C212.7,93.1 206.9,96.0 205.4,96.6 C205.1,102.4 203.0,107.8 198.3,112.5 C181.9,128.9 168.3,122.5 157.7,114.1 C157.9,116.9 156.7,120.9 152.7,124.9 L141.0,136.5 C139.8,137.7 141.6,141.9 141.8,141.8 Z" fill="currentColor" class="octo-body"></path></svg><style>.github-corner:hover .octo-arm{animation:octocat-wave 560ms ease-in-out}@keyframes octocat-wave{0%,100%{transform:rotate(0)}20%,60%{transform:rotate(-25deg)}40%,80%{transform:rotate(10deg)}}@media (max-width:500px){.github-corner:hover .octo-arm{animation:none}.github-corner .octo-arm{animation:octocat-wave 560ms ease-in-out}}</style></a>
<style>
    td 
    {
	background-color: lightgrey;
    }
</style>
<?php
include 'db.inc.php';

//TODO: manage json export
//TODO: manage type from DDB
//TODO: manage type from Flarm
//TODO: manage timezone (query + display)
//date_default_timezone_set("Zulu");

if (isset($_GET['airfield'])) { 
	if (isset($_GET['d'])) {
		$date = DateTime::createFromFormat('!Y-m-d', $_GET['d']);
	} else {
		$date = new DateTime();
	}
	$date->setTime(0,0,0); // Be sure to be at begining of the day.
	$requestedDate = clone($date);
	$date->add(new DateInterval('P1D')); // Add a day to be at the end of the day
	$nextDate = clone($date);
	$date->sub(new DateInterval('P2D'));
	$previousDate = clone($date);

	$handle = $link->prepare('SELECT 
		id,
		aircraftId,
		aircraftReg,
		aircraftCN,
		aircraftTracked,
		aircraftIdentified,
		UNIX_TIMESTAMP(takeoffTimestamp) as takeoffTimestamp,
		takeoffAirfield,
		UNIX_TIMESTAMP(landingTimestamp) as landingTimestamp,
		landingAirfield
		FROM flightlog
		WHERE ( `takeoffAirfield` = :airfield OR `landingAirfield` = :airfield )
		 AND ((`takeoffTimestamp` > FROM_UNIXTIME(:dateMin)
		 AND `takeoffTimestamp` < FROM_UNIXTIME(:dateMax))
		 OR (`landingTimestamp` > FROM_UNIXTIME(:dateMin)
                 AND `landingTimestamp` < FROM_UNIXTIME(:dateMax)))
		 ORDER BY `id`');
	$handle->bindValue(':airfield', $_GET['airfield']);
	$handle->bindValue(':dateMin', $requestedDate->getTimestamp());
	$handle->bindValue(':dateMax', $nextDate->getTimestamp());
	$handle->execute();
	$data = $handle->fetchAll();
	echo '<center><h1>Flightlog for '.htmlspecialchars($_GET['airfield']).'</h1>';
	echo '<h2>';
	echo '<a href="?airfield='.$_GET['airfield'].'&d='.$previousDate->format('Y-m-d').'">&lt;</a> ';
	echo $requestedDate->format('Y-m-d');
	echo ' <a href="?airfield='.$_GET['airfield'].'&d='.$nextDate->format('Y-m-d').'">&gt;</a>';
	echo '</h2><br>';
	echo '<table style="text-align:center;"><tr>'.
		"<th>ID</th>".
		"<th>Reg</th>".
		"<th>CN</th>".
		"<th>Takeoff time</th>".
		"<th>Takeoff airfield</th>".
		"<th>Landing time</th>".
		"<th>Landing airfield</th>".
		"</tr>";
	foreach ($data as $row) {
		if ($row['aircraftTracked'] == "Y" && $row['aircraftIdentified'] == "Y") {
			echo "<tr><td>".
				$row['aircraftId']."</td><td>".
				$row['aircraftReg']."</td><td>".
				$row['aircraftCN']."</td><td>";
		} else { // Notrack in DDB or No identified in DDB => hidden
			echo "<tr><td colspan=\"3\"><font size=\"-1\">(hidden)</font></td><td>";
		}
		if ($row['takeoffTimestamp'] > 0)
			echo date('r',$row['takeoffTimestamp']);
		echo 
			"</td><td>".
			$row['takeoffAirfield']."</td><td>";
		if ($row['landingTimestamp'] > 0)
			echo date('r',$row['landingTimestamp']);
		echo 
			"</td><td>".
			$row['landingAirfield']."</td></tr>\n";
	}
	echo "</table></center>";
} else {
	$handle = $link->prepare('SELECT `takeoffAirfield` as `airfield` FROM `flightlog` UNION SELECT `landingAirfield` FROM `flightlog` ORDER BY `airfield`');
	$handle->execute();
	$data = $handle->fetchAll();
?>
<center>
Please select an airfield:<br>
<?php
	foreach ($data as $row) {
		echo '<a href="?airfield='.htmlspecialchars($row['airfield']).'">'.htmlspecialchars($row['airfield']).'</a><br>';
	}
?>
</center>
<?php
}
?>
</body>
</html>
