

	var addr1 = "ws://" + document.location.host + "/body";
	conn_body = new WebSocket(addr1);
	conn_body.onmessage = function (evt) {
		var messages = evt.data.split('@');
		var type = messages[0];
		var id = messages[1];
		var item = document.getElementById(id);
		if (type === "TEXT") {
			item.innerHTML = messages[2];
		}
		else if (type === "COLOR") {
			var color = messages[2];
			item.style.color = color;
		}
		else if (type === "BCOLOR") {
			var color = messages[2];
			item.style.backgroundColor  = color;
		}
		else if (type === "FONTSIZE") {
			var fsize = messages[2];
			item.style.fontSize  = fsize;
		}
		else if (type === "FONTFAMILY") {
			var font = messages[2];
			item.style.fontFamily  = font;
		}
		else if (type === "MODALSHOW") {
			var modal = new bootstrap.Modal(item); 
			modal.show();
		}
		else if (type === "ENABLE") {
			var enable = messages[2];
			if (enable === "ENABLE") {
				item.disabled = false;
			}
			else  {
				item.disabled = true;
			}
		}
		else if (type === "IMAGE") {
			src = 'static/' + messages[2]; 
			item.src = src
		}
		
	};
	
	function bt1_func() {
		xhr = new XMLHttpRequest();
		xhr.open("POST", "/bt1", true);
		xhr.send();
	}
	