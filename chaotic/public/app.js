var Chaotic = {
    enableForm: function() {
	$('form input').removeAttr('disabled');
	$('#spin').hide();
	$('#save').show();
    },

    disableForm: function() {
	$('form input').attr('disabled', true);
	$('#save').hide();
	$('#spin').show();
    },

    populate: function(policy) {
	["Delay", "DelayP", "FailureP"].forEach(function(s) {
	    $('form [name='+s+']').val(policy[s]);
	});
	this.enableForm();
    },

    load: function() {
	var that = this;
	$.getJSON("policy",
		  function(data) { that.populate(data); });
    },

    save: function() {
	// Serializing the form into JSON data
	var v = function(name) { return $('form input[name='+name+']').val() };
	var data = {
	    Delay: v('Delay'),
	    DelayP: parseFloat(v('DelayP')),
	    FailureP: parseFloat(v('FailureP'))
	}

	this.disableForm();
	var that = this;
	$.post("policy", JSON.stringify(data),
	       function(data) { that.populate(data); });
    },

    listen: function() {
	var that = this;
	$('form').on('submit', function(event) {
	    event.preventDefault();
	    that.save();
	});
    },

    init: function() {
	this.load();
	this.listen();
    }
};

$(function() { Chaotic.init(); });
