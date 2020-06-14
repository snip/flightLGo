<?php
include 'db.inc.php';
if (isset($_POST['takeoffTimestamp'])) { // Takeoff call
	$handle = $link->prepare('INSERT INTO flightlog (`lastIP`,`aircraftId`,`aircraftReg`,`aircraftCN`,`aircraftTracked`,`aircraftIdentified`,`takeoffTimestamp`,`takeoffAirfield`) 
		VALUES(:lastIP,:aircraftId,:aircraftReg,:aircraftCN,:aircraftTracked,:aircraftIdentified,FROM_UNIXTIME(:takeoffTimestamp),:takeoffAirfield)');
	$handle->bindValue(':lastIP', $_SERVER['REMOTE_ADDR']);
	$handle->bindValue(':aircraftId', $_POST['aircraftId']);
	$handle->bindValue(':aircraftReg', $_POST['aircraftReg']);
	$handle->bindValue(':aircraftCN', $_POST['aircraftCN']);
	$handle->bindValue(':aircraftTracked', $_POST['aircraftTracked']);
	$handle->bindValue(':aircraftIdentified', $_POST['aircraftIdentified']);
	$handle->bindValue(':takeoffTimestamp', $_POST['takeoffTimestamp']);
	$handle->bindValue(':takeoffAirfield', $_POST['takeoffAirfield']);
	$handle->execute();
} elseif (isset($_POST['landingTimestamp'])) { // Landing call
	// Add a new raw if not already detected as flying
	// Otherwise update existing raw
	// TODO: Check landing time is after takoff time (really required?)
  	$handle = $link->prepare('
INSERT INTO flightlog (`lastIP`,`aircraftId`,`aircraftReg`,`aircraftCN`,`aircraftTracked`,`aircraftIdentified`,`landingTimestamp`,`landingAirfield`)
	SELECT :lastIP,:aircraftId,:aircraftReg,:aircraftCN,:aircraftTracked,:aircraftIdentified,FROM_UNIXTIME(:landingTimestamp),:landingAirfield FROM dual 
	WHERE NOT EXISTS (SELECT 1 FROM flightlog WHERE `aircraftId`=:aircraftId AND `landingTimestamp` IS NULL AND `takeoffTimestamp` > DATE_SUB(NOW(), INTERVAL 1 DAY) LIMIT 1);
UPDATE flightlog SET `lastIP`=:lastIP,`landingTimestamp`=FROM_UNIXTIME(:landingTimestamp),`landingAirfield`=:landingAirfield
	WHERE `aircraftId`=:aircraftId AND `landingTimestamp` IS NULL AND `takeoffTimestamp` > DATE_SUB(NOW(), INTERVAL 1 DAY) ORDER BY `takeoffTimestamp` DESC LIMIT 1
	');
	$handle->bindValue(':lastIP', $_SERVER['REMOTE_ADDR']);
	$handle->bindValue(':aircraftId', $_POST['aircraftId']);
	$handle->bindValue(':aircraftReg', $_POST['aircraftReg']);
	$handle->bindValue(':aircraftCN', $_POST['aircraftCN']);
	$handle->bindValue(':aircraftTracked', $_POST['aircraftTracked']);
	$handle->bindValue(':aircraftIdentified', $_POST['aircraftIdentified']);
	$handle->bindValue(':landingTimestamp', $_POST['landingTimestamp']);
	$handle->bindValue(':landingAirfield', $_POST['landingAirfield']);
	$handle->execute();
}
?>
