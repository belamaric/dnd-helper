
function getMonsters(name, cb) {
  $.get({
    url: "/monsters",
    data: {
      search: name
    },
    success: cb
    error: func() { alert("error") }
  });
}

function searchMonsters(name) {
	
}

function showMonster(name) {
	alert("Showing " + name)
	getMonsters(name, function(result) {
		alert(result["monsters"][0]["Name"])
	}
	)
}

