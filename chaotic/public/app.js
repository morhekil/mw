var Chaotic = {}

Chaotic.Form = {
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

Chaotic.Log = {
    canvas: null,
    width: 1000,
    height: 100,
    barw: 4,
    x: null,
    y: null,
    live: true,
    paused: false,

    timeFn: function(d) { return d.Time },
    indexFn: function(d) { return d.Index },
    
    color: function(d) {
	if (d.Failed) {
	    return 'red';
	} else if (d.Delayed) {
	    return 'orange';
	} else {
	    return '#888';
	}
    },

    graph: function(data) {
	// $('#log').html(data);
	var maxn = d3.max(data, this.indexFn);
	console.log('maxn', maxn);
	this.x.domain([maxn - Math.floor(this.width / this.barw), maxn])
	var miny = d3.min(data, this.timeFn);
	this.y.domain([miny > 0 ? miny : 1, d3.max(data, this.timeFn)]);

	var that = this;
	var bars = this.canvas.selectAll('rect')
	    .data(data, this.indexFn)

	bars.enter()
	    .append('svg:rect')
	    .attr('x', function(d) { return that.x(d.Index + 1) })
	    .attr('y', function(d) { return that.height - that.y(d.Time) })
	    .attr('width', this.barw - 1)
	    .attr('height', function(d) { return that.y(d.Time); })
	    .style('fill', this.color)
	    .append('svg:title')
	    .text(function(d) { return JSON.stringify(d) });

	bars.transition()
	    .attr('x', function(d) { return that.x(d.Index) })
	    .attr('y', function(d) { return that.height - that.y(d.Time) })
	    .attr('height', function(d) { return that.y(d.Time); })
	
    },

    load: function() {
	var that = this;
	d3.json("log", function(err, data) {
	    if (err) { return $('#log').addClass('error').html(err); }
	    that.graph(data);
	});

    },

    init: function() {
	var that = this;
	this.canvas = d3.select('#log').append('svg:svg')
	    .attr('width', this.width)
	    .attr('height', this.height)
	    .on('mouseover', function() { that.paused = true; })
	    .on('mouseout', function() { that.paused = false; });

	this.y = d3.scale.log().clamp(true).range([0, this.height])
	this.x = d3.scale.linear().range([0, this.width - this.barw])
	that.load();
	window.setInterval(function() {
	    if (that.live && !that.paused) { that.load(); }
	}, 500);
    }
}

$(function() {
    Chaotic.Form.init();
    Chaotic.Log.init();
});
