var Chaotic = function() {
    var populate = function(policy) {
	    $('form [name=Delay]').val(policy.Delay);
	    $('form [name=DelayP]').val(policy.DelayP);
	    $('form input').removeAttr('disabled');
    }

    var load = function() {
	$.getJSON("policy", populate);
    };

    var save = function() {
	// Serializing the form into JSON data
	var v = function(name) { return $('form input[name='+name+']').val() };
	var data = {
	    Delay: v('Delay'),
	    DelayP: parseFloat(v('DelayP'))
	}

	$.post("policy", JSON.stringify(data), populate);
    };

    var listen = function() {
	$('form').on('submit', function(event) {
	    event.preventDefault();
	    save();
	});
    }
	    
    load();
    listen();
};

$(function() {
    Chaotic();
});
