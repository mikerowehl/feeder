function markRead(t) {
	console.log("Going to make a request for " + t.getAttribute('id'));
	t.innerHTML = "Requesting...";
	var xhr = new XMLHttpRequest();
	xhr.open("POST", '/markread', true);
	xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
	xhr.onreadystatechange = function() {
		if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
			t.innerHTML = "Processed!!!";
		}
	}
	xhr.send("guid=" + t.getAttribute('id'));
}

window.onload = function() {
	var x = document.getElementsByClassName("marker");
	var i;
	for (i=0; i < x.length; i++) {
		x[i].onclick = function handler(e) { markRead(e.target) }
	}
}
