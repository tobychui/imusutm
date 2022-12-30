<?php
$tokens = [];


if (!isset($_GET['token'])){
	header("HTTP/1.1 401 Unauthorized");
	echo "Token not given";
    exit;
}else{
	$reqToken = $_GET['token'];
	if (in_array($reqToken, $tokens)){
		$json = file_get_contents("http://localhost:8089/");
		header('Content-Type: application/json; charset=utf-8');
		echo $json;
	}else{
		header("HTTP/1.1 401 Unauthorized");
		echo "Invalid token";
		exit;
	}
}


?>