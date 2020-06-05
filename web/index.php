<html>
<body>
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
		 AND (`takeoffTimestamp` > FROM_UNIXTIME(:dateMin)
		 AND `takeoffTimestamp` < FROM_UNIXTIME(:dateMax))
		 OR (`landingTimestamp` > FROM_UNIXTIME(:dateMin)
                 AND `landingTimestamp` < FROM_UNIXTIME(:dateMax))
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
