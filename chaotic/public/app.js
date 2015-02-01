var Chaotic = {}

//
// Policy form operations and management
//
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

//
// Log data visualisation
//
Chaotic.Log = {
    canvas: null,
    bars: null,
    width: 1000,
    height: 100,
    barw: 4,
    x: null,
    y: null,
    live: true,
    paused: false,
    focused: null,
    delayIntv: null,

    timeFn: function(d) { return d.Time },
    indexFn: function(d) { return d.Index },
    
    color: function(d) {
	var unfocused = (this.focused && this.focused != d.Index) ||
	    (this.filter && !d.Text.match(this.filter));
	if (unfocused) {
	    return '#ddd';
	} else if (d.Failed) {
	    return 'red';
	} else if (d.Delayed) {
	    return 'orange';
	} else {
	    return '#888';
	}
    },

    domains: function(data) {
	var maxn = d3.max(data, this.indexFn);
	this.x.domain([maxn - Math.floor(this.width / this.barw), maxn])
	var miny = d3.min(data, this.timeFn);
	this.y.domain([miny > 0 ? miny : 1, d3.max(data, this.timeFn)]);
    },

    graph: function(data) {
	this.domains(data);

	this.bars = this.canvas.selectAll('rect')
	    .data(data, this.indexFn)
	this.enter();
	this.transition();
	this.bars.exit().remove();
    },

    recolorise: function(delay) {
	var that = this;
	if (this.delayIntv) { window.clearInterval(this.delayIntv); }

	window.setTimeout(function() {
	    that.bars
		.style('fill', function(d) { return that.color(d) });
	}, delay);
    },

    // convert nanoseconds to a more human-readable representation
    nsToHuman: function(ns) {
	if (ns < 1e3) {
	    return ns + "ns";
	} else if (ns < 1e6) {
	    return Math.round(ns / 1e3) + "&mu;s"
	} else if (ns < 1e9) {
	    return Math.round(ns / 1e6) + "ms"
	} else {
	    return Math.round(ns / 1e9 * 10) / 10 + "s"
	}
    },

    annotate: function(d) {
	$('#log .annotate').html(d.Text + "<br />" + this.nsToHuman(d.Time));
	this.focused = d.Index;
	this.recolorise(1);
    },

    clearAnnotate: function() {
	$('#log .annotate').html('');
	this.focused = null;
	this.recolorise(1000);
    },

    enter: function() {
	var that = this;
	this.bars.enter()
	    .append('svg:rect')
	    .attr('x', function(d) { return that.x(d.Index + 1) })
	    .attr('y', function(d) { return that.height - that.y(d.Time) })
	    .attr('width', this.barw - 1)
	    .attr('height', function(d) { return that.y(d.Time); })
	    .style('fill', function(d) { return that.color(d) })
	    .on('mouseover', function(d) { that.annotate(d); })
	    .on('mouseout', function() { that.clearAnnotate(); })
    },

    transition: function() {
	var that = this;
	this.bars.transition()
	    .attr('x', function(d) { return that.x(d.Index) })
	    .attr('y', function(d) { return that.height - that.y(d.Time) })
	    .attr('height', function(d) { return that.y(d.Time); })
	
    },

    load: function() {
	var that = this;
	d3.json("log", function(err, data) {
	    if (err) { return $('#log').addClass('error'); }
	    $('#log').removeClass('error');
	    that.graph(data);
	});

    },

    scaleCanvas: function() {
	this.width = $('.container')[0].getBoundingClientRect().width;
	this.canvas.attr('width', this.width);
	this.x = d3.scale.linear().range([0, this.width - this.barw])
	this.y = d3.scale.log().clamp(true).range([0, this.height])
    },

    clear: function(e) {
	e.preventDefault();
	var btn = $('#clear');
	btn.removeClass('button-primary').attr('disabled', true);
	$.post('log', function() {
	    btn.addClass('button-primary').removeAttr('disabled');
	});
    },

    initCanvas: function() {
	var that = this;
	this.canvas = d3.select('#log').append('svg:svg')
	    .attr('height', this.height)
	    .on('mouseover', function() { that.paused = true; })
	    .on('mouseout', function() { that.paused = false; });
	this.scaleCanvas();
    },

    initControls: function() {
	var that = this;
	$('#clear').on('click', function(e) { that.clear(e); });
	$('#live').on('change', function() {
	    that.live = $(this).prop('checked');
	});
	$('#filter').on('keyup', function() {
	    that.filter = new RegExp($(this).val(), 'i');
	    that.recolorise(300);
	});
    },
    
    init: function() {
	var that = this;
	this.initCanvas();
	this.initControls();

	that.load();
	window.setInterval(function() {
	    if (that.live && !that.paused) { that.load(); }
	}, 500);
	window.onresize = function() { that.scaleCanvas(); }
    }
}

$(function() {
    Chaotic.Form.init();
    Chaotic.Log.init();
});
