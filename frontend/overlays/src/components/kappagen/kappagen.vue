<!-- eslint-disable no-undef -->
<!-- eslint-disable no-case-declarations -->
<!-- eslint-disable no-prototype-builtins -->
<script setup>
import './kappagen.css';
import { ref } from 'vue';

const stupidStylesRef = ref(null);

const cfg = {
	display: {
		styles: [
			'Still',        // No movement
			'StraightLine', // Gentle movement in a random straight line
			'Rise',         // Slowly rise to top while wobbling back and forth
			'Bounce',       // Fall from the top at an angle and bounce along the bottom (Windows Solitaire style)
			'Speed',        // Zoom across the screen
			'Drop',         // Get stuck at the top and tumble down (no fade/zoom in, only out)
			'Crazy',        // Bounce off the walls [Layout Shifts - Requires offset-anchor / offset-position directives]
			'Confetti',     // Fall like confetti                   (no zoom, no fade in, only fade out)
			'Throw',        // Toss at the middle and tumble down   (no fade/zoom in, only out)
			'TheCube',       // Rotate a 3D cube of an emote         (no zoom, only fade)
		],
		kappa: {
			count: 150,
			styles: {
				'Burst': {},       // Expand from a center point
				'Fireworks': {},   // Burst out from a single emote (no fade/zoom; small emotes)
				'Spiral': {},      // Burst out from a single emote in a spiral (no fade/zoom; small emotes)
				'Pyramid': {},     // Build a pyramid        (no fade/zoom; specific-size emotes)
				'SmallPyramid': {}, // Build a small pyramid  (no fade/zoom; small emotes)
				'Fountain': {},    // Spout from a fountain  (no fade/zoom)
				'Stampede': {},    // Stampede of emotes     (no fade/zoom)
				'Confetti': {},    // Fall like confetti     (no zoom, no fade in, only fade out; small emotes)
				'Conga': {},       // Start a conga line     (no fade/zoom)
				'TheCube': {       // Rotate a 3D cube of an emote (no zoom, only fade)
					size: 8 / 10,
					center: true,
					rotations: 5,
					faces: false,
				},
				'Text': {          // Show a message         (no fade/zoom; specific-size emotes)
					message: ['HYPE!'],
					time: 3,
				},
			},
			conga: {
				contagious: false,
				time: 5,
				avoidMiddle: false,
			},
		},
		statuses: true,
		extended: {
			useFFZ: true,
			useBTTV: true,
			use7TV: true,
			useZWE: true,
			fillZWE: false,
		},
	},
	emote: {
		time: 5,
		max: 0,
		queue: 0,
		size: {
			ratio: {
				normal: 1 / 12,
				small: 1 / 24,
			},
			min: 16,
			max: 256,
			variation: false,
		},
		cube: {
			rotations: 5,
		},
		in: {
			fade: true,
			zoom: true, /* Layout Shifts - Requires independent scale directive */
		},
		out: {
			fade: true,
			zoom: true, /* Layout Shifts - Requires independent scale directive */
		},
	},
};

/* fractions (or decimal percentages) of the emote time configuration value */
const timing = {
	display: {
		'Still': {
			time: 1,
		},
		'StraightLine': {
			time: 1,
		},
		'Rise': {
			origin: {      /* percentages of the screen height */
				min: 0.8,
				max: 1.1,
			},
			time: 1,
			wiggle: {      /* percentages of the above time percentage */
				delay: {
					min: 0,
					max: 3 / 25,
				},
				min: 2 / 5,
				max: 1,
			},
		},
		'Bounce': {
			origin: {      /* percentages of the screen height */
				min: 0,
				max: 0.2,
			},
			time: 1,
			velocity: {    /* pixels per 300th of display time */
				h: {
					min: 3,
					max: 9,
				},
				v: {
					min: 4,
					max: 7,
				},
				loss: 0.3,     /* velocity percentage lost per bounce */
			},
			gravity: 1,     /* pixels added to vertical velocity per increment */
		},
		'Speed': {
			origin: {      /* percentages of the screen height */
				min: 0.3,
				max: 0.7,
			},
			time: 1,
			delay: 0.1,
		},
		'Drop': {
			time: 1,
		},
		'Crazy': {
			time: 1,
			distance: 7000,  /* max pixels to travel */
			squash: {
				width: 2,     /* squashed wide dimension */
				height: 0.7,  /* squashed tall dimension */
				time: 1 / 50,
			},
		},
		'Confetti': {
			time: 1,
		},
		'Throw': {
			time: 1,
			twist: 7 / 50,
			toss: 1 / 5,
			drop: 4 / 5,
			dest: {
				h: {          /* percentages of the screen width */
					min: 0.3,
					max: 0.7,
				},
				v: {          /* percentages of the screen height */
					min: 0.3,
					max: 0.7,
				},
			},
		},
		'TheCube': {
			time: 1,
		},
		'Fountain': {
			time: 1 / 2,
		},
	},
	kappa: {
		'Rise': {
			time: 2,
		},
		'Speed': {
			time: 2,
		},
		'Crazy': {
			time: 2,
		},
		'Burst': {
			time: 1.5,
			top: {         /* top and bottom margin of the origin point */
				min: 1 / 4,
				max: 3 / 4,
			},             /* left and right margin of the origin point */
			left: {
				min: 1 / 4,
				max: 3 / 4,
			},
		},
		'Fireworks': {
			time: 1,
			origin: {      /* origin point(s) of the firework's rocket */
				x: [1 / 2],
				y: [1],
			},
			dest: {        /* destination point(s) of the firework's rocket */
				x: [1 / 4, 1 / 2, 3 / 4],
				y: [1 / 3],
			},
			speed: {
				rocket: 2 / 5,  /* speed of rocket */
				burst: 1 / 50,   /* speed of initial burst */
			},
			quantity: {    /* number of emotes per burst */
				small: 1 / 8,
				medium: 3 / 4,
				large: 1 / 8,
			},
			radius: {      /* firework burst radii */
				base: 2 / 3,    /* screen's smaller dimension */
				small: 1 / 3,
				medium: 2 / 3,
				large: 1,
			},
			spread: 12,    /* how much more frequently to pause during medium burst */
			delays: {      /* pause between bursts */
				small: 2 / 25,
				large: 1 / 10,
			},
		},
		'Spiral': {
			time: 1 / 2,
			bulk: 8,       /* max number of emotes to send in bulk (>1 can end up looking chunked) */
			vectors: {     /* number of emote vectors per circle */
				min: 40,
				max: 60,
			},
		},
		'Pyramid': {
			time: 1,
			show: {
				total: 0.8,   /* percentage of time to spend showing the pyramid */
				min: 75,       /* minimum animation speed per block (in ms) */
			},
			pause: 0.2,
			hide: 0.01,
		},
		'SmallPyramid': {
			time: 1,
			show: {
				total: 0.8,   /* percentage of time to spend showing the pyramid */
				min: 100,      /* minimum animation speed per block (in ms) */
			},
			pause: 0.2,
			hide: 0.01,
		},
		'Fountain': {
			time: 1.5,
			top: {         /* peak of the fountain, as a percent of the screen height */
				min: 3 / 20,
				max: 2 / 5,
			},             /* left and right margin of the origin point */
			left: {
				min: 1 / 3,
				max: 2 / 3,
			},
		},
		'Stampede': {
			time: 1,
			speed: 2 / 5,     /* travel time across the screen for each emote */
			maxdensity: 6,  /* maximum emotes to show at once */
			top: {
				min: 0.5,      /* top of stampede relative to top of screen, in emote heights */
				max: 0.5,       /* bottom of stampede relative to bottom of screen, in emote heights */
			},
			height: 3,      /* height of stampede in emote heights */
			bunch: {
				'1': {
					min: 1,
					max: 5,
				},
				'2': 8,        /* this number minus the value of 1 */
				'4': {
					min: 0,
					max: 3,
				},
			},
			pause: {
				'1': 4 / 5,
				'2': 2 / 5,
			},
			smallSleep: {
				min: 90,
				max: 100,
			},
		},
		'Confetti': {
			time: 1,
		},
		'Conga': {
			time: {
				show: 2,
				hide: 2,
			},
			size: 5 / 3,      /* height of animation space for row in emote heights */
			height: 5 / 6,    /* height of each row of dancers in emote heights (padding) */
			avoidMiddle: 6,  /* rows to use when avoiding the middle (half top, half bottom; please use even numbers) */
		},
		'TheCube': {
			time: 1,
		},
		'Text': {
			time: 1,
			show: {
				total: 0.8,   /* percentage of time to spend showing the pyramid */
				min: 75,       /* minimum animation speed per block (in ms) */
			},
			hide: 0.01,
		},
	},
};

/* potentially alterable arrays */

// list of default images to use if your channel has no emotes
const bareList = [
	{ url: 'https://cdn.7tv.app/emote/6548b7074789656a7be787e1/4x.webp' },
];

// distribution of emotes for Pyramid and SmallPyramid kappagen
const pyramidDist = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1];

// treat these BTTV emotes as zero-width (overlapping) emotes
const bttvZWE = [
	'567b5b520e984428652809b6', //SoSnowy
	'5849c9a4f52be01a7ee5f79d', //IceCold
	'58487cc6f52be01a7ee5f205', //SantaHat
	'5849c9c8f52be01a7ee5f79e', //TopHat
	'567b5dc00e984428652809bd', //ReinDeer
	'567b5c080e984428652809ba', //CandyCane
	'5e76d399d6581c3724c0f0b8', //cvMask
	'5e76d338d6581c3724c0f0b2',  //cvHazmat
];

// distribution of emotes for letters in message kappagens
const alnumDist = {
	'A': [
		[0, 1, 1, 1, 1, 1, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 1, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 0, 1, 0, 1, 0, 0],
		[0, 1, 1, 1, 1, 1, 0, 0, 0],
	],
	'a': [
		[0, 0, 1, 1, 1, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 1, 1, 1, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
	],
	'B': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 1, 1, 0, 1, 1, 0, 0],
	],
	'b': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 0, 1, 1, 1, 0, 0, 0, 0],
	],
	'C': [
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 1, 0, 0, 0, 1, 0, 0],
	],
	'c': [
		[0, 0, 1, 1, 1, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
	],
	'D': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
	],
	'd': [
		[0, 0, 1, 1, 1, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
	],
	'E': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
	],
	'e': [
		[0, 0, 1, 1, 1, 0, 0, 0, 0],
		[0, 1, 0, 1, 0, 1, 0, 0, 0],
		[0, 1, 0, 1, 0, 1, 0, 0, 0],
		[0, 1, 0, 1, 0, 1, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
	],
	'F': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
	],
	'f': [
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 1, 1, 1, 1, 1, 1, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 0, 1, 0, 1, 0, 0],
	],
	'G': [
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 1, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 1, 1, 1, 0, 1, 0, 0],
	],
	'g': [
		[1, 0, 0, 1, 1, 0, 0, 0, 0],
		[1, 0, 1, 0, 0, 1, 0, 0, 0],
		[1, 0, 1, 0, 0, 1, 0, 0, 0],
		[0, 1, 1, 1, 1, 0, 0, 0, 0],
	],
	'H': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
	],
	'h': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 1, 1, 1, 0, 0, 0, 0, 0],
	],
	'I': [
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
	],
	'i': [
		[0, 1, 0, 0, 1, 0, 0, 0, 0],
		[0, 1, 1, 1, 1, 0, 1, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
	],
	'J': [
		[0, 0, 1, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
	],
	'j': [
		[1, 0, 0, 0, 0, 0, 0, 0, 0],
		[1, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 1, 1, 1, 1, 0, 1, 0, 0],
	],
	'K': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 1, 0, 1, 0, 0, 0],
		[0, 0, 1, 0, 0, 0, 1, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
	],
	'k': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 1, 0, 0, 0, 0, 0],
		[0, 0, 1, 0, 1, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
	],
	'L': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
	],
	'l': [
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
	],
	'M': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 0, 0, 1, 0, 0],
		[0, 0, 0, 0, 0, 1, 0, 0, 0],
		[0, 0, 0, 0, 0, 0, 1, 0, 0],
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
	],
	'm': [
		[0, 1, 1, 1, 1, 1, 0, 0, 0],
		[0, 0, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 1, 1, 1, 0, 0, 0, 0],
		[0, 0, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 1, 1, 1, 0, 0, 0, 0],
	],
	'N': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 0, 0, 1, 0, 0],
		[0, 0, 0, 1, 1, 1, 0, 0, 0],
		[0, 0, 1, 0, 0, 0, 0, 0, 0],
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
	],
	'n': [
		[0, 1, 1, 1, 1, 1, 0, 0, 0],
		[0, 0, 0, 0, 0, 1, 0, 0, 0],
		[0, 0, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 1, 1, 1, 0, 0, 0, 0],
	],
	'O': [
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
	],
	'o': [
		[0, 0, 1, 1, 1, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 0, 1, 1, 1, 0, 0, 0, 0],
	],
	'P': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 0, 0, 1, 1, 0, 0],
	],
	'p': [
		[1, 1, 1, 1, 1, 0, 0, 0, 0],
		[0, 0, 1, 0, 0, 1, 0, 0, 0],
		[0, 0, 1, 0, 0, 1, 0, 0, 0],
		[0, 0, 0, 1, 1, 0, 0, 0, 0],
	],
	'Q': [
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 1, 0, 0, 0, 0, 1, 0],
		[1, 1, 0, 1, 1, 1, 1, 0, 0],
	],
	'q': [
		[0, 0, 0, 1, 1, 0, 0, 0, 0],
		[0, 0, 1, 0, 0, 1, 0, 0, 0],
		[0, 0, 1, 0, 0, 1, 0, 0, 0],
		[1, 1, 1, 1, 1, 0, 0, 0, 0],
	],
	'R': [
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 1, 1, 0, 0, 1, 0],
		[0, 0, 1, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 1, 1, 0, 0],
	],
	'r': [
		[0, 1, 1, 1, 1, 1, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 0, 0, 1, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
	],
	'S': [
		[0, 0, 1, 0, 0, 1, 1, 0, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 1, 1, 0, 0, 1, 0, 0],
	],
	's': [
		[0, 1, 0, 0, 1, 0, 0, 0, 0],
		[0, 1, 0, 1, 0, 1, 0, 0, 0],
		[0, 1, 0, 1, 0, 1, 0, 0, 0],
		[0, 1, 0, 1, 0, 1, 0, 0, 0],
		[0, 0, 1, 0, 0, 1, 0, 0, 0],
	],
	'T': [
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
	],
	't': [
		[0, 0, 0, 0, 0, 1, 0, 0, 0],
		[0, 0, 0, 0, 0, 1, 0, 0, 0],
		[0, 0, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 0, 0, 0, 0, 1, 0, 0, 0],
	],
	'U': [
		[0, 0, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 0, 1, 1, 1, 1, 1, 1, 0],
	],
	'u': [
		[0, 0, 1, 1, 1, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 1, 1, 1, 1, 0, 0, 0],
	],
	'V': [
		[0, 0, 0, 0, 1, 1, 1, 1, 0],
		[0, 0, 1, 1, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 0, 1, 1, 0, 0, 0, 0, 0],
		[0, 0, 0, 0, 1, 1, 1, 1, 0],
	],
	'v': [
		[0, 0, 0, 0, 1, 1, 0, 0, 0],
		[0, 0, 1, 1, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 0, 1, 1, 0, 0, 0, 0, 0],
		[0, 0, 0, 0, 1, 1, 0, 0, 0],
	],
	'W': [
		[0, 0, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 0, 1, 1, 1, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 0, 1, 1, 1, 1, 1, 1, 0],
	],
	'w': [
		[0, 0, 1, 1, 1, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 0, 1, 1, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 0, 1, 1, 1, 1, 0, 0, 0],
	],
	'X': [
		[0, 1, 1, 0, 0, 0, 1, 1, 0],
		[0, 0, 0, 1, 0, 1, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 1, 0, 1, 0, 0, 0],
		[0, 1, 1, 0, 0, 0, 1, 1, 0],
	],
	'x': [
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 0, 1, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 1, 0, 0, 0, 0, 0],
		[0, 0, 1, 0, 1, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
	],
	'Y': [
		[0, 0, 0, 0, 0, 0, 1, 1, 0],
		[0, 0, 0, 0, 1, 1, 0, 0, 0],
		[0, 1, 1, 1, 0, 0, 0, 0, 0],
		[0, 0, 0, 0, 1, 1, 0, 0, 0],
		[0, 0, 0, 0, 0, 0, 1, 1, 0],
	],
	'y': [
		[1, 0, 0, 1, 1, 1, 0, 0, 0],
		[1, 0, 1, 0, 0, 0, 0, 0, 0],
		[1, 0, 1, 0, 0, 0, 0, 0, 0],
		[0, 1, 1, 1, 1, 1, 0, 0, 0],
	],
	'Z': [
		[0, 1, 1, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 1, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 1, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 1, 1, 0],
	],
	'z': [
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[0, 1, 1, 0, 0, 1, 0, 0, 0],
		[0, 1, 0, 1, 0, 1, 0, 0, 0],
		[0, 1, 0, 0, 1, 1, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
	],
	'1': [
		[0, 1, 0, 0, 0, 0, 1, 0, 0],
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
	],
	'2': [
		[0, 1, 1, 0, 0, 0, 1, 0, 0],
		[0, 1, 0, 1, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 1, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 1, 0, 0],
	],
	'3': [
		[0, 0, 1, 0, 0, 0, 1, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 1, 1, 0, 1, 1, 0, 0],
	],
	'4': [
		[0, 0, 0, 1, 1, 0, 0, 0, 0],
		[0, 0, 0, 1, 0, 1, 1, 1, 0],
		[0, 0, 0, 1, 0, 0, 0, 0, 0],
		[0, 1, 1, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 1, 0, 0, 0, 0, 0],
	],
	'5': [
		[0, 0, 1, 0, 0, 1, 1, 1, 0],
		[0, 1, 0, 0, 0, 1, 0, 1, 0],
		[0, 1, 0, 0, 0, 1, 0, 1, 0],
		[0, 0, 1, 1, 1, 0, 0, 1, 0],
	],
	'6': [
		[0, 0, 1, 1, 1, 1, 0, 0, 0],
		[0, 1, 0, 0, 1, 0, 1, 0, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 1, 1, 0, 0, 0, 0, 0],
	],
	'7': [
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 1, 1, 0, 0, 0, 1, 0],
		[0, 0, 0, 0, 1, 1, 0, 1, 0],
		[0, 0, 0, 0, 0, 0, 1, 1, 0],
	],
	'8': [
		[0, 0, 1, 1, 0, 1, 1, 0, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 1, 1, 0, 1, 1, 0, 0],
	],
	'9': [
		[0, 0, 0, 0, 0, 1, 1, 0, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 1, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 1, 1, 1, 1, 0, 0],
	],
	'0': [
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
		[0, 1, 1, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 1, 1, 1, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 1, 1, 0],
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
	],
	'>': [
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 1, 0, 0, 0, 1, 0, 0],
		[0, 0, 0, 1, 0, 1, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
	],
	'<': [
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 1, 0, 1, 0, 0, 0],
		[0, 0, 1, 0, 0, 0, 1, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
	],
	':': [
		[0, 0, 1, 0, 0, 1, 0, 0, 0],
	],
	'.': [
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
	],
	',': [
		[1, 0, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
	],
	'\'': [
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
		[0, 0, 0, 0, 0, 0, 0, 0, 1],
	],
	'-': [
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
	],
	'_': [
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 0, 0, 0, 0],
	],
	'+': [
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
	],
	'=': [
		[0, 0, 0, 1, 0, 1, 0, 0, 0],
		[0, 0, 0, 1, 0, 1, 0, 0, 0],
		[0, 0, 0, 1, 0, 1, 0, 0, 0],
	],
	'!': [
		[0, 0, 0, 0, 0, 1, 1, 0, 0],
		[0, 1, 0, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 0, 1, 1, 0, 0],
	],
	'@': [
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 1, 0, 1, 0, 1, 0],
		[0, 1, 0, 1, 0, 1, 0, 1, 0],
		[0, 1, 0, 0, 1, 1, 1, 0, 0],
	],
	'#': [
		[0, 0, 0, 0, 1, 0, 1, 0, 0],
		[0, 0, 0, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 1, 0, 1, 0, 0],
		[0, 0, 0, 1, 1, 1, 1, 1, 0],
		[0, 0, 0, 0, 1, 0, 1, 0, 0],
	],
	'$': [
		[0, 0, 1, 0, 0, 1, 1, 0, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[1, 1, 1, 1, 1, 1, 1, 1, 1],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 1, 1, 0, 0, 1, 0, 0],
	],
	'\u00a2': [
		[0, 0, 1, 1, 1, 0, 0, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
		[1, 1, 1, 1, 1, 1, 1, 0, 0],
		[0, 1, 0, 0, 0, 1, 0, 0, 0],
	],
	'\u20ac': [
		[0, 0, 0, 1, 0, 1, 0, 0, 0],
		[0, 0, 1, 1, 1, 1, 1, 0, 0],
		[0, 1, 0, 1, 0, 1, 0, 1, 0],
		[0, 1, 0, 1, 0, 1, 0, 1, 0],
		[0, 1, 0, 1, 0, 1, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 0, 1, 0],
	],
	'\u00a3': [
		[0, 1, 0, 0, 1, 0, 1, 0, 0],
		[0, 1, 1, 1, 1, 1, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 1, 0, 0, 1, 0],
		[0, 1, 0, 0, 0, 0, 1, 0, 0],
	],
	'\u00a5': [
		[0, 0, 0, 0, 1, 0, 1, 1, 0],
		[0, 0, 1, 0, 1, 1, 0, 0, 0],
		[0, 1, 1, 1, 1, 0, 0, 0, 0],
		[0, 0, 1, 0, 1, 1, 0, 0, 0],
		[0, 0, 0, 0, 1, 0, 1, 1, 0],
	],
	'%': [
		[0, 1, 0, 0, 0, 0, 1, 1, 0],
		[0, 0, 1, 1, 0, 0, 1, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 0, 0],
		[0, 1, 1, 0, 0, 1, 1, 0, 0],
		[0, 1, 1, 0, 0, 0, 0, 1, 0],
	],
	'?': [
		[0, 0, 0, 0, 0, 1, 1, 0, 0],
		[0, 0, 0, 0, 0, 0, 0, 1, 0],
		[0, 1, 0, 1, 1, 0, 0, 1, 0],
		[0, 0, 0, 0, 1, 0, 0, 1, 0],
		[0, 0, 0, 0, 0, 1, 1, 0, 0],
	],
};

//////////////////////////////////////////////////////////////////////////////
// don't mess with things below this line without knowing what you're doing //
//////////////////////////////////////////////////////////////////////////////

const display = function () {
	let _eActive = 0;
	let _iTitanic = 0;
	const _cRadius = Math.PI * 2;
	const _tAnim = {
		fade: {
			in: 8,
			out: 8,
		},
		zoom: {
			in: 17,
			out: 8,
		},
	};

	const $emote = function () {
		const _toShow = [];

		let _tEmote = false;

		const $list = function () {
			function $Still(eInf, sW, sH, eH, canV = true, tInit = 0) {
				if (tInit === 0)
					tInit = new Date().getTime();
				if (_iTitanic > tInit)
					return;
				let variationSize = 1;
				if (canV && cfg.emote.size.variation !== false) {
					if (typeof cfg.emote.size.variation === 'number') {
						const chances = [];
						chances.push(0.5);
						chances.push(2);
						for (let i = 0; i < cfg.emote.size.variation; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
						const chances = [];
						chances.push(...cfg.emote.size.variation.range);
						for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					}
				}
				eH = Math.ceil(eH * variationSize);
				let eW = eH;
				if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
					eW = eInf.width / eInf.height * eH;
				const h = shared.random(sW - eW);
				const v = shared.random(sH - eH);
				let s = 'top: ' + v + 'px;';
				s += ' left: ' + h + 'px;';
				s += ' --emote-height: ' + eH + 'px;';
				s += ' --emote-width: ' + eW + 'px;';
				const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.Still.time);
				s += _styleEmote([], [], [], [], [], [], cfg.emote.in.fade, cfg.emote.in.zoom, cfg.emote.out.fade, cfg.emote.out.zoom, tMS);
				_addEmoteToDoc(tInit, eInf.url, variationSize, { style: s }, false, {
					space: false,
					time: tMS,
				});
				if (eInf.hasOwnProperty('zwe')) {
					for (let i = 0, l = eInf.zwe.length; i < l; i++) {
						_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, { style: s }, false, {
							space: false,
							time: tMS,
						});
					}
				}
			}

			function $StraightLine(eInf, sW, sH, eH, x = false, y = false, canV = true, tInit = 0) {
				if (tInit === 0)
					tInit = new Date().getTime();
				if (_iTitanic > tInit)
					return;
				let variationSize = 1;
				if (canV && cfg.emote.size.variation !== false) {
					if (typeof cfg.emote.size.variation === 'number') {
						const chances = [];
						chances.push(0.5);
						chances.push(2);
						for (let i = 0; i < cfg.emote.size.variation; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
						const chances = [];
						chances.push(...cfg.emote.size.variation.range);
						for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					}
				}
				eH = Math.ceil(eH * variationSize);
				let eW = eH;
				if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
					eW = eInf.width / eInf.height * eH;
				const eHh = Math.ceil(eH / 2);
				const eWh = Math.ceil(eW / 2);
				let h = x;
				if (h === false)
					h = shared.random(sW) - eWh;
				let v = y;
				if (v === false)
					v = shared.random(sH) - eHh;
				const r = Math.min(sW, sH) * (shared.random() + 1);
				let th = shared.random() * _cRadius;
				if (!x && !y) {
					const nH = eH * -1;
					const nW = eW * -1;
					while (!_safePoints(h, v, th, r, nW, nH, sW, sH)) {
						th = shared.random() * _cRadius;
					}
				}
				const hD = Math.floor(h + r * Math.cos(th));
				const vD = Math.floor(v + r * Math.sin(th));
				let s = '--emote-height: ' + eH + 'px;';
				s += ' --emote-width: ' + eW + 'px;';
				const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.StraightLine.time);
				s += ' transform: translate(' + h + 'px, ' + v + 'px);';
				s += _styleEmote([], [], [], [], [], [], cfg.emote.in.fade, cfg.emote.in.zoom, cfg.emote.out.fade, cfg.emote.out.zoom, tMS);
				_addEmoteToDoc(tInit, eInf.url, variationSize, {
					style: s,
					classes: ['etStraightLine'],
				}, false, { time: tMS }, { x: hD, y: vD });
				if (eInf.hasOwnProperty('zwe')) {
					for (let i = 0, l = eInf.zwe.length; i < l; i++) {
						_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, {
							style: s,
							classes: ['etStraightLine'],
						}, false, { time: tMS }, { x: hD, y: vD });
					}
				}
			}

			function $Rise(eInf, sW, sH, eH, canV = true, tInit = 0) {
				if (tInit === 0)
					tInit = new Date().getTime();
				if (_iTitanic > tInit)
					return;
				let variationSize = 1;
				if (canV && cfg.emote.size.variation !== false) {
					if (typeof cfg.emote.size.variation === 'number') {
						const chances = [];
						chances.push(0.5);
						chances.push(2);
						for (let i = 0; i < cfg.emote.size.variation; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
						const chances = [];
						chances.push(...cfg.emote.size.variation.range);
						for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					}
				}
				eH = Math.ceil(eH * variationSize);
				let eW = eH;
				if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
					eW = eInf.width / eInf.height * eH;
				const eWh = Math.ceil(eW / 2);
				const h = shared.random(sW) - eWh;
				const v = Math.floor(sH * _rndFromRange(timing.display.Rise.origin));
				let s = 'left: ' + h + 'px;';
				s += ' --emote-height: ' + eH + 'px;';
				s += ' --emote-width: ' + eW + 'px;';
				if (cfg.emote.out.fade || cfg.emote.out.zoom)
					s += ' offset-path: path("M 0 ' + v + ' L 0 ' + Math.floor(v * 0.05) + '") ;';
				else
					s += ' offset-path: path("M 0 ' + v + ' L 0 -' + eH + '") ;';
				const aNames = [];
				const aDelays = [];
				const aDurs = [];
				const aTimings = [];
				const aFills = [];
				const aIters = [];
				if (shared.random(2) === 0)
					aNames.push('wiggleL');
				else
					aNames.push('wiggleR');
				const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.Rise.time);
				const d = Math.floor(tMS * _rndFromRange(timing.display.Rise.wiggle.delay));
				aDelays.push(d + 'ms');
				const w = Math.floor(tMS * _rndFromRange(timing.display.Rise.wiggle));
				aDurs.push(w + 'ms');
				aTimings.push('ease-in-out');
				aFills.push('both');
				aIters.push('infinite');
				aNames.push('offsetPath');
				aDelays.push('0s');
				aDurs.push(tMS + 'ms');
				aTimings.push('linear');
				aFills.push('forwards');
				aIters.push('1');
				s += _styleEmote(aNames, aDelays, aDurs, aTimings, aFills, aIters, cfg.emote.in.fade, cfg.emote.in.zoom, cfg.emote.out.fade, cfg.emote.out.zoom, tMS);
				_addEmoteToDoc(tInit, eInf.url, variationSize, { style: s }, false, {
					space: false,
					time: tMS,
				});
				if (eInf.hasOwnProperty('zwe')) {
					for (let i = 0, l = eInf.zwe.length; i < l; i++) {
						_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, { style: s }, false, {
							space: false,
							time: tMS,
						});
					}
				}
			}

			const $Bounce = function () {
				function $c_Bounce(eInf, sW, sH, eH, canV = true, tInit = 0) {
					if (tInit === 0)
						tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					let variationSize = 1;
					if (canV && cfg.emote.size.variation !== false) {
						if (typeof cfg.emote.size.variation === 'number') {
							const chances = [];
							chances.push(0.5);
							chances.push(2);
							for (let i = 0; i < cfg.emote.size.variation; i++) {
								chances.push(1);
							}
							variationSize = chances[shared.random(chances.length)];
						} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
							const chances = [];
							chances.push(...cfg.emote.size.variation.range);
							for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
								chances.push(1);
							}
							variationSize = chances[shared.random(chances.length)];
						}
					}
					eH = Math.ceil(eH * variationSize);
					let eW = eH;
					if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
						eW = eInf.width / eInf.height * eH;
					const eWh = Math.ceil(eW / 2);
					const sWm = Math.ceil(sW / 2);
					const h = Math.floor(shared.random(sW) - eWh);
					const v = Math.floor(sH * _rndFromRange(timing.display.Bounce.origin));
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.Bounce.time);
					const vMS = (tMS / 300 / (16 + 2 / 3));
					let velH = _rndFromRange(timing.display.Bounce.velocity.h);
					const velV = _rndFromRange(timing.display.Bounce.velocity.v);
					if (h + eWh > sWm)
						velH = -1 * velH;
					let s = '--emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eW + 'px;';
					s += _styleEmote([], [], [], [], [], [], cfg.emote.in.fade, cfg.emote.in.zoom, cfg.emote.out.fade, cfg.emote.out.zoom, tMS);
					s += ' transform: translate(' + h + 'px, ' + v + 'px);';
					const bX = h;
					const bY = v;
					const iArr = [];
					iArr.push(_addEmoteToDoc(tInit, eInf.url, variationSize, { style: s }, true, { time: tMS }));
					if (eInf.hasOwnProperty('zwe')) {
						for (let i = 0, l = eInf.zwe.length; i < l; i++) {
							iArr.push(_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, { style: s }, true, { time: tMS }));
						}
					}
					window.requestAnimationFrame(function (ts) {
						_tLoop(tInit, iArr, bX, bY, velH, velV, vMS, sH, eH, ts, ts);
					});
				}

				function _tLoop(tInit, iArr, bX, bY, velH, velV, vMS, sH, eH, myT, ts) {
					if (_iTitanic > tInit)
						return;
					if (iArr[0].parentElement === null)
						return;
					let steps = 1;
					if (myT === 0)
						myT = ts;
					else {
						steps = Math.max(1, Math.floor((ts - myT) / 16));
						myT = ts;
					}
					for (let i = 0; i < steps; i++) {
						bX += velH / vMS;
						bY += velV / vMS;
						velV += timing.display.Bounce.gravity / vMS;
						const sB = sH - eH;
						if (bY >= sB) {
							bY = sB;
							velV *= -1 * (1 - timing.display.Bounce.velocity.loss);
							velV = Math.floor(velV);
						}
					}
					for (let i = 0, l = iArr.length; i < l; i++) {
						iArr[i].style.transform = 'translate(' + bX + 'px, ' + bY + 'px)';
					}
					window.requestAnimationFrame(function (fTS) {
						_tLoop(tInit, iArr, bX, bY, velH, velV, vMS, sH, eH, myT, fTS);
					});
				}

				return $c_Bounce;
			}();

			function $Speed(eInf, sW, sH, eH, canV = true, tInit = 0) {
				if (tInit === 0)
					tInit = new Date().getTime();
				if (_iTitanic > tInit)
					return;
				let variationSize = 1;
				if (canV && cfg.emote.size.variation !== false) {
					if (typeof cfg.emote.size.variation === 'number') {
						const chances = [];
						chances.push(0.5);
						chances.push(2);
						for (let i = 0; i < cfg.emote.size.variation; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
						const chances = [];
						chances.push(...cfg.emote.size.variation.range);
						for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					}
				}
				eH = Math.ceil(eH * variationSize);
				let eW = eH;
				if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
					eW = eInf.width / eInf.height * eH;
				const eWh = Math.ceil(eW / 2);
				const sWm = Math.ceil(sW / 2);
				const h = shared.random(sW) - eWh;
				const v = Math.floor(sH * _rndFromRange(timing.display.Speed.origin));
				let s = 'top: ' + v + 'px;';
				s += ' left: ' + h + 'px;';
				s += ' --emote-height: ' + eH + 'px;';
				s += ' --emote-width: ' + eW + 'px;';
				const aNames = [];
				const aDelays = [];
				const aDurs = [];
				const aTimings = [];
				const aFills = [];
				const aIters = [];
				const dsO = {};
				if (h + eWh > sWm) {
					dsO.origin = 'right';
					aNames.push('speedL');
				} else {
					dsO.origin = 'left';
					aNames.push('speedR');
				}
				const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.Speed.time);
				const d = Math.floor(tMS * timing.display.Speed.delay);
				aDelays.push(d + 'ms');
				aDurs.push((tMS - d) + 'ms');
				aTimings.push('ease-in');
				aFills.push('forwards');
				aIters.push('1');
				s += _styleEmote(aNames, aDelays, aDurs, aTimings, aFills, aIters, cfg.emote.in.fade, cfg.emote.in.zoom, cfg.emote.out.fade, cfg.emote.out.zoom, tMS);
				_addEmoteToDoc(tInit, eInf.url, variationSize, {
					style: s,
					dataset: dsO,
				}, false, { time: tMS });
				if (eInf.hasOwnProperty('zwe')) {
					for (let i = 0, l = eInf.zwe.length; i < l; i++) {
						_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, {
							style: s,
							dataset: dsO,
						}, false, { time: tMS });
					}
				}
			}

			function $Drop(eInf, sW, sH, eH, canV = true, tInit = 0) {
				if (tInit === 0)
					tInit = new Date().getTime();
				if (_iTitanic > tInit)
					return;
				let variationSize = 1;
				if (canV && cfg.emote.size.variation !== false) {
					if (typeof cfg.emote.size.variation === 'number') {
						const chances = [];
						chances.push(0.5);
						chances.push(2);
						for (let i = 0; i < cfg.emote.size.variation; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
						const chances = [];
						chances.push(...cfg.emote.size.variation.range);
						for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					}
				}
				eH = Math.ceil(eH * variationSize);
				let eW = eH;
				if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
					eW = eInf.width / eInf.height * eH;
				const eWh = Math.ceil(eW / 2);
				const h = shared.random(sW) - eWh;
				let s = 'left: ' + h + 'px;';
				s += ' --emote-height: ' + eH + 'px;';
				s += ' --emote-width: ' + eW + 'px;';
				const aNames = [];
				const aDelays = [];
				const aDurs = [];
				const aTimings = [];
				const aFills = [];
				const aIters = [];
				const dsO = {};
				if (shared.random(2) === 0) {
					dsO.origin = 'topleft';
					aNames.push('dropL');
				} else {
					dsO.origin = 'topright';
					aNames.push('dropR');
				}
				aDelays.push('0s');
				const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.Drop.time);
				aDurs.push(tMS + 'ms');
				aTimings.push('ease-in');
				aFills.push('forwards');
				aIters.push('1');
				s += _styleEmote(aNames, aDelays, aDurs, aTimings, aFills, aIters, false, false, cfg.emote.out.fade, cfg.emote.out.zoom, tMS);
				_addEmoteToDoc(tInit, eInf.url, variationSize, {
					style: s,
					dataset: dsO,
				}, false, { space: false, time: tMS });
				if (eInf.hasOwnProperty('zwe')) {
					for (let i = 0, l = eInf.zwe.length; i < l; i++) {
						_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, {
							style: s,
							dataset: dsO,
						}, false, { space: false, time: tMS });
					}
				}
			}

			const $Crazy = function () {
				/* LAYOUT SHIFTS
             * =============
             * squashes via scale
             * offset-path requires support for offset-anchor/offset-position
             * due to transform-origin changes during squash
             */

				function $c_Crazy(eInf, sW, sH, eH, canV = true, tInit = 0) {
					if (tInit === 0)
						tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					let variationSize = 1;
					if (canV && cfg.emote.size.variation !== false) {
						if (typeof cfg.emote.size.variation === 'number') {
							const chances = [];
							chances.push(0.5);
							chances.push(2);
							for (let i = 0; i < cfg.emote.size.variation; i++) {
								chances.push(1);
							}
							variationSize = chances[shared.random(chances.length)];
						} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
							const chances = [];
							chances.push(...cfg.emote.size.variation.range);
							for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
								chances.push(1);
							}
							variationSize = chances[shared.random(chances.length)];
						}
					}
					eH = Math.ceil(eH * variationSize);
					let eW = eH;
					if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
						eW = eInf.width / eInf.height * eH;
					const sR = sW - eW;
					const sB = sH - eH;
					const h = shared.random(sR - 5) + 10;
					const v = shared.random(sB - 5) + 10;
					let s = 'top: ' + v + 'px;';
					s += ' left: ' + h + 'px;';
					s += ' --emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eW + 'px;';
					const dests = [];
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.Crazy.time);
					const rate = Math.sqrt(timing.display.Crazy.distance ** 2 / 2) / tMS;
					const traj = { x: 0, y: 0 };
					while (traj.x === 0 && traj.y === 0) {
						traj.x = shared.random() * (rate * 2) - rate;
						traj.y = shared.random() * (rate * 2) - rate;
					}
					const pos = { x: h, y: v, t: 0 };
					let lastT = 0;
					let bCt = 0;
					const sqTime = Math.floor(tMS * timing.display.Crazy.squash.time * 2);
					while (bCt * sqTime + pos.t < tMS) {
						pos.x += traj.x;
						pos.y += traj.y;
						pos.t += 1;
						let wall = false;
						if (pos.x <= 0) {
							pos.x = 0;
							traj.x *= -1;
							wall = 1;
						} else if (pos.x >= sR) {
							pos.x = sR;
							traj.x *= -1;
							wall = 3;
						}
						if (pos.y <= 0) {
							pos.y = 0;
							traj.y *= -1;
							wall = 2;
						} else if (pos.y >= sB) {
							pos.y = sB;
							traj.y *= -1;
							wall = 4;
						}
						if (wall !== false) {
							bCt++;
							dests.push({
								x: Math.floor(pos.x),
								y: Math.floor(pos.y),
								t: pos.t - lastT,
								w: wall,
							});
							lastT = pos.t;
						}
					}
					dests.push({ x: pos.x, y: pos.y, t: pos.t - lastT, w: 0 });
					s += _styleEmote([], [], [], [], [], [], cfg.emote.in.fade, cfg.emote.in.zoom, cfg.emote.out.fade, cfg.emote.out.zoom, tMS);
					const iArr = [];
					iArr.push(_addEmoteToDoc(tInit, eInf.url, variationSize, { style: s }, true, {
						space: false,
						time: tMS,
					}));
					if (eInf.hasOwnProperty('zwe')) {
						for (let i = 0, l = eInf.zwe.length; i < l; i++) {
							iArr.push(_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, { style: s }, true, {
								space: false,
								time: tMS,
							}));
						}
					}
					const d = 0;
					const lA = iArr.length;
					const lD = dests.length;
					shared.doNextFrame(_tLoop, tInit, lA, iArr, lD, dests, d);
				}

				function _tLoop(tInit, lA, iArr, lD, dests, d) {
					if (_iTitanic > tInit)
						return;
					if (d >= lD)
						return;
					const squashT = Math.floor(cfg.emote.time * 1000 * timing.display.Crazy.time * timing.display.Crazy.squash.time);
					for (let i = 0; i < lA; i++) {
						iArr[i].dataset.origin = 'center';
						iArr[i].dataset.squash = 'no';
						iArr[i].style.top = dests[d].y + 'px';
						iArr[i].style.left = dests[d].x + 'px';
						iArr[i].style.transition = 'top ' + dests[d].t + 'ms linear, left ' + dests[d].t + 'ms linear, transform ' + squashT + 'ms linear';
					}
					d++;
					window.setTimeout(_tSquash, dests[d - 1].t, tInit, lA, iArr, lD, dests, d);
				}

				function _tSquash(tInit, lA, iArr, lD, dests, d) {
					if (_iTitanic > tInit)
						return;
					const squashT = Math.floor(cfg.emote.time * 1000 * timing.display.Crazy.time * timing.display.Crazy.squash.time);
					for (let i = 0; i < lA; i++) {
						switch (dests[d - 1].w) {
							case 1:
								iArr[i].dataset.origin = 'left';
								iArr[i].dataset.squash = 'horizontal';
								break;
							case 2:
								iArr[i].dataset.origin = 'top';
								iArr[i].dataset.squash = 'vertical';
								break;
							case 3:
								iArr[i].dataset.origin = 'right';
								iArr[i].dataset.squash = 'horizontal';
								break;
							case 4:
								iArr[i].dataset.origin = 'bottom';
								iArr[i].dataset.squash = 'vertical';
								break;
						}
					}
					window.setTimeout(_tUnsquash, squashT, tInit, lA, iArr, lD, dests, d);
				}

				function _tUnsquash(tInit, lA, iArr, lD, dests, d) {
					if (_iTitanic > tInit)
						return;
					const squashT = Math.floor(cfg.emote.time * 1000 * timing.display.Crazy.time * timing.display.Crazy.squash.time);
					for (let i = 0; i < lA; i++) {
						iArr[i].dataset.squash = 'no';
					}
					window.setTimeout(_tLoop, squashT, tInit, lA, iArr, lD, dests, d);
				}

				return $c_Crazy;
			}();

			function $Confetti(eInf, sW, sH, eH, canV = true, tInit = 0) {
				if (tInit === 0)
					tInit = new Date().getTime();
				if (_iTitanic > tInit)
					return;
				let variationSize = 1;
				if (canV && cfg.emote.size.variation !== false) {
					if (typeof cfg.emote.size.variation === 'number') {
						const chances = [];
						chances.push(0.5);
						chances.push(2);
						for (let i = 0; i < cfg.emote.size.variation; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
						const chances = [];
						chances.push(...cfg.emote.size.variation.range);
						for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					}
				}
				eH = Math.ceil(eH * variationSize);
				let eW = eH;
				if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
					eW = eInf.width / eInf.height * eH;
				const eWh = Math.ceil(eW / 2);
				const h = shared.random(sW) - eWh;
				let s = 'left: ' + h + 'px;';
				s += ' --emote-height: ' + eH + 'px;';
				s += ' --emote-width: ' + eW + 'px;';
				const aNames = [];
				const aDelays = [];
				const aDurs = [];
				const aTimings = [];
				const aFills = [];
				const aIters = [];
				switch (shared.random(3)) {
					case 0:
						aNames.push('confettiA');
						break;
					case 1:
						aNames.push('confettiB');
						break;
					case 2:
						aNames.push('confettiC');
						break;
				}
				aDelays.push('0s');
				const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.Confetti.time);
				aDurs.push(tMS + 'ms');
				aTimings.push('linear');
				aFills.push('forwards');
				aIters.push('1');
				s += _styleEmote(aNames, aDelays, aDurs, aTimings, aFills, aIters, false, false, cfg.emote.out.fade, false, tMS);
				_addEmoteToDoc(tInit, eInf.url, variationSize, { style: s }, false, {
					space: false,
					time: tMS,
				});
				if (eInf.hasOwnProperty('zwe')) {
					for (let i = 0, l = eInf.zwe.length; i < l; i++) {
						_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, { style: s }, false, {
							space: false,
							time: tMS,
						});
					}
				}
			}

			const $Throw = function () {
				function $c_Throw(eInf, sW, sH, eH, canV = true, tInit = 0) {
					if (tInit === 0)
						tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					let variationSize = 1;
					if (canV && cfg.emote.size.variation !== false) {
						if (typeof cfg.emote.size.variation === 'number') {
							const chances = [];
							chances.push(0.5);
							chances.push(2);
							for (let i = 0; i < cfg.emote.size.variation; i++) {
								chances.push(1);
							}
							variationSize = chances[shared.random(chances.length)];
						} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
							const chances = [];
							chances.push(...cfg.emote.size.variation.range);
							for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
								chances.push(1);
							}
							variationSize = chances[shared.random(chances.length)];
						}
					}
					eH = Math.ceil(eH * variationSize);
					let eW = eH;
					if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
						eW = eInf.width / eInf.height * eH;
					const sR = sW - eW;
					const sB = sH - eH;
					const h = shared.random(2) === 0 ? eW * -1 : sW;
					const v = shared.random(sH + eH) - eH;
					const hD = Math.floor(sR * _rndFromRange(timing.display.Throw.dest.h));
					const vD = Math.floor(sB * _rndFromRange(timing.display.Throw.dest.v));
					const dH = shared.random() * eH;
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.Throw.time);
					const t2 = Math.floor(tMS * timing.display.Throw.toss);
					const t3 = Math.floor(tMS * timing.display.Throw.drop);
					let s = '--emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eW + 'px;';
					s += ' transform: translate(' + h + 'px, ' + v + 'px);';
					let s2 = '--emote-height: ' + eH + 'px;';
					s2 += ' --emote-width: ' + eW + 'px;';
					let r = '360';
					if (h > 0)
						r = '-360';
					s2 += ' transform: translate(' + hD + 'px, ' + (sH - dH) + 'px) rotate(' + r + 'deg);';
					const aNames = [];
					const aDelays = [];
					const aDurs = [];
					const aTimings = [];
					const aFills = [];
					const aIters = [];
					if (cfg.emote.out.fade) {
						const fOut = _tAnim.fade.out / 100;
						const t3F = t3 * fOut;
						aNames.push('fadeOut');
						aDelays.push(Math.floor(t3 - t3F) + 'ms');
						aDurs.push(Math.floor(t3F) + 'ms');
						aTimings.push('ease-out');
						aFills.push('forwards');
						aIters.push('1');
					} else {
						aNames.push('noFadeOut');
						aDelays.push(t3 - 50 + 'ms');
						aDurs.push('50ms');
						aTimings.push('ease-out');
						aFills.push('forwards');
						aIters.push('1');
					}
					if (cfg.emote.out.zoom) {
						const zOut = _tAnim.zoom.out / 100;
						const t3Z = t3 * zOut;
						aNames.push('zoomOut');
						aDelays.push(Math.floor(t3 - t3Z) + 'ms');
						aDurs.push(Math.floor(t3Z) + 'ms');
						aTimings.push('linear');
						aFills.push('forwards');
						aIters.push('1');
					}
					s2 += _styleEmoteString(aNames, aDelays, aDurs, aTimings, aFills, aIters, tMS);
					const iArr = [];
					iArr.push(_addEmoteToDoc(tInit, eInf.url, variationSize, {
						style: s,
						classes: ['etThrowTwist'],
					}, true, { space: false, time: tMS }));
					if (eInf.hasOwnProperty('zwe')) {
						for (let i = 0, l = eInf.zwe.length; i < l; i++) {
							iArr.push(_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, {
								style: s,
								classes: ['etThrowTwist'],
							}, true, { space: false, time: tMS }));
						}
					}
					shared.doNextFrame(_tMove, tInit, iArr, hD, vD);
					window.setTimeout(_tDrop, t2, tInit, iArr, s2);
				}

				function _tMove(tInit, iArr, hD, vD) {
					if (_iTitanic > tInit)
						return;
					for (let i = 0, l = iArr.length; i < l; i++) {
						iArr[i].style.transform = 'translate(' + hD + 'px, ' + vD + 'px)';
					}
				}

				function _tDrop(tInit, iArr, s2) {
					if (_iTitanic > tInit)
						return;
					for (let i = 0, l = iArr.length; i < l; i++) {
						iArr[i].classList.replace('etThrowTwist', 'etThrowDrop');
						iArr[i].setAttribute('style', s2);
					}
				}

				return $c_Throw;
			}();

			const $TheCube = function () {
				function $c_TheCube(eInf, sW, sH, eH, canV = true, tInit = 0) {
					if (tInit === 0)
						tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					const scene = document.createElement('div');
					scene.setAttribute('class', 'scene fit cube');
					let variationSize = 1;
					if (canV && cfg.emote.size.variation !== false) {
						if (typeof cfg.emote.size.variation === 'number') {
							const chances = [];
							chances.push(0.5);
							chances.push(2);
							for (let i = 0; i < cfg.emote.size.variation; i++) {
								chances.push(1);
							}
							variationSize = chances[shared.random(chances.length)];
						} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
							const chances = [];
							chances.push(...cfg.emote.size.variation.range);
							for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
								chances.push(1);
							}
							variationSize = chances[shared.random(chances.length)];
						}
					}
					eH = Math.ceil(eH * variationSize);
					const eHh = Math.ceil(eH / 2);
					const nHh = eHh * -1;
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.TheCube.time);
					const cube = document.createElement('div');
					cube.setAttribute('class', 'cube');
					cube.setAttribute('style', 'transform: translateZ(' + nHh + 'px);');
					if (!eInf.hasOwnProperty('zwe') || eInf.zwe.length === 0) {
						const cubeF = document.createElement('img');
						_setImgSrc(cubeF, eInf.url);
						cubeF.dataset.face = 'front';
						cube.appendChild(cubeF);
						const cubeB = document.createElement('img');
						_setImgSrc(cubeB, eInf.url);
						cubeB.dataset.face = 'back';
						cube.appendChild(cubeB);
						const cubeR = document.createElement('img');
						_setImgSrc(cubeR, eInf.url);
						cubeR.dataset.face = 'right';
						cube.appendChild(cubeR);
						const cubeL = document.createElement('img');
						_setImgSrc(cubeL, eInf.url);
						cubeL.dataset.face = 'left';
						cube.appendChild(cubeL);
						const cubeT = document.createElement('img');
						_setImgSrc(cubeT, eInf.url);
						cubeT.dataset.face = 'top';
						cube.appendChild(cubeT);
						const cubeU = document.createElement('img');
						_setImgSrc(cubeU, eInf.url);
						cubeU.dataset.face = 'bottom';
						cube.appendChild(cubeU);
					} else {
						const lZ = eInf.zwe.length;
						const cubeF = document.createElement('div');
						cubeF.dataset.face = 'front';
						const pctF = document.createElement('img');
						_setImgSrc(pctF, eInf.url);
						cubeF.appendChild(pctF);
						for (let i = 0; i < lZ; i++) {
							const pctZ = document.createElement('img');
							_setImgSrc(pctZ, eInf.zwe[i].url);
							cubeF.appendChild(pctZ);
						}
						cube.appendChild(cubeF);
						const cubeB = document.createElement('div');
						cubeB.dataset.face = 'back';
						const pctB = document.createElement('img');
						_setImgSrc(pctB, eInf.url);
						cubeB.appendChild(pctB);
						for (let i = 0; i < lZ; i++) {
							const pctZ = document.createElement('img');
							_setImgSrc(pctZ, eInf.zwe[i].url);
							cubeB.appendChild(pctZ);
						}
						cube.appendChild(cubeB);
						const cubeR = document.createElement('div');
						cubeR.dataset.face = 'right';
						const pctR = document.createElement('img');
						_setImgSrc(pctR, eInf.url);
						cubeR.appendChild(pctR);
						for (let i = 0; i < lZ; i++) {
							const pctZ = document.createElement('img');
							_setImgSrc(pctZ, eInf.zwe[i].url);
							cubeR.appendChild(pctZ);
						}
						cube.appendChild(cubeR);
						const cubeL = document.createElement('div');
						cubeL.dataset.face = 'left';
						const pctL = document.createElement('img');
						_setImgSrc(pctL, eInf.url);
						cubeL.appendChild(pctL);
						for (let i = 0; i < lZ; i++) {
							const pctZ = document.createElement('img');
							_setImgSrc(pctZ, eInf.zwe[i].url);
							cubeL.appendChild(pctZ);
						}
						cube.appendChild(cubeL);
						const cubeT = document.createElement('div');
						cubeT.dataset.face = 'top';
						const pctT = document.createElement('img');
						_setImgSrc(pctT, eInf.url);
						cubeT.appendChild(pctT);
						for (let i = 0; i < lZ; i++) {
							const pctZ = document.createElement('img');
							_setImgSrc(pctZ, eInf.zwe[i].url);
							cubeT.appendChild(pctZ);
						}
						cube.appendChild(cubeT);
						const cubeU = document.createElement('div');
						cubeU.dataset.face = 'bottom';
						const pctU = document.createElement('img');
						_setImgSrc(pctU, eInf.url);
						cubeU.appendChild(pctU);
						for (let i = 0; i < lZ; i++) {
							const pctZ = document.createElement('img');
							_setImgSrc(pctZ, eInf.zwe[i].url);
							cubeU.appendChild(pctZ);
						}
						cube.appendChild(cubeU);
					}
					scene.appendChild(cube);
					const h = shared.random(sW) - eHh;
					const v = shared.random(sH) - eHh;
					const r = Math.min(sW, sH) * (shared.random() + 1);
					let th = shared.random() * _cRadius;
					const nH = eH * -1;
					while (!_safePoints(h, v, th, r, nH, nH, sW, sH)) {
						th = shared.random() * _cRadius;
					}
					const hD = Math.floor(h + r * Math.cos(th));
					const vD = Math.floor(v + r * Math.sin(th));
					let s = '--emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eH + 'px;';
					s += ' --cube-depth: ' + eHh + 'px;';
					s += ' perspective: ' + eH * 3 + 'px;';
					s += ' transform: translate(' + h + 'px, ' + v + 'px);';
					s += _styleEmote([], [], [], [], [], [], cfg.emote.in.fade, false, cfg.emote.out.fade, false, tMS);
					scene.setAttribute('style', s);
					_eActive += 6;
					document.body.appendChild(scene);
					_gc.hook(scene, true, 6, tMS);
					shared.doNextFrame(_tMove, tInit, cube, scene, hD, vD, eH);
				}

				function _tMove(tInit, cube, scene, hD, vD, eH) {
					if (_iTitanic > tInit)
						return;
					const nHh = Math.ceil(eH / 2) * -1;
					let rX = 0;
					let rY = 0;
					while (Math.abs(rX) + Math.abs(rY) < 45) {
						rX = (360 - shared.random() * 720) * cfg.emote.cube.rotations;
						rY = (360 - shared.random() * 720) * cfg.emote.cube.rotations;
					}
					cube.style.transform = 'translateZ(' + nHh + 'px) rotateX(' + rX + 'deg) rotateY(' + rY + 'deg)';
					scene.style.transform = 'translate(' + hD + 'px, ' + vD + 'px)';
				}

				return $c_TheCube;
			}();

			function $Fountain(eInf, sW, sH, eH, fX, fY, canV = true, tInit = 0) {
				if (tInit === 0)
					tInit = new Date().getTime();
				if (_iTitanic > tInit)
					return;
				const tMS = Math.floor(cfg.emote.time * 1000 * timing.display.Fountain.time);
				let variationSize = 1;
				if (canV && cfg.emote.size.variation !== false) {
					if (typeof cfg.emote.size.variation === 'number') {
						const chances = [];
						chances.push(0.5);
						chances.push(2);
						for (let i = 0; i < cfg.emote.size.variation; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					} else if (typeof cfg.emote.size.variation === 'object' && cfg.emote.size.variation.hasOwnProperty('chance') && cfg.emote.size.variation.hasOwnProperty('range') && Array.isArray(cfg.emote.size.variation.range)) {
						const chances = [];
						chances.push(...cfg.emote.size.variation.range);
						for (let i = 0; i < cfg.emote.size.variation.chance; i++) {
							chances.push(1);
						}
						variationSize = chances[shared.random(chances.length)];
					}
				}
				eH = Math.ceil(eH * variationSize);
				let eW = eH;
				if (eInf.hasOwnProperty('width') && eInf.hasOwnProperty('height'))
					eW = eInf.width / eInf.height * eH;
				const sR = sW - eW;
				const sB = sH - eH;
				let h = fX;
				if (h === false)
					h = Math.floor(shared.random() * (sR * 0.33) + sR * 0.33);
				let hD;
				if (shared.random(2) === 0)
					hD = h - shared.random(sR * 0.2);
				else
					hD = h + shared.random(sR * 0.2);
				let s = '--emote-height: ' + eH + 'px;';
				s += ' --emote-width: ' + eW + 'px;';
				s += ' transform: translateX(' + h + 'px);';
				s += ' offset-path: path("M 0 ' + sH + ' L 0 ' + Math.floor(fY * sH + shared.random(sB / 2)) + ' L 0 ' + (sH + eH) + '");';
				const aNames = [];
				const aDelays = [];
				const aDurs = [];
				const aTimings = [];
				const aFills = [];
				const aIters = [];
				aNames.push('offsetPath');
				aDelays.push('0s');
				aDurs.push(tMS + 'ms');
				aTimings.push('cubic-bezier(0, 0.9, 1, 0.15)');
				aFills.push('forwards');
				aIters.push('1');
				s += _styleEmoteString(aNames, aDelays, aDurs, aTimings, aFills, aIters, tMS);
				_addEmoteToDoc(tInit, eInf.url, variationSize, {
					style: s,
					classes: ['etFountain'],
				}, false, { time: tMS, space: false }, { x: hD });
				if (eInf.hasOwnProperty('zwe')) {
					for (let i = 0, l = eInf.zwe.length; i < l; i++) {
						_addEmoteToDoc(tInit, eInf.zwe[i].url, variationSize, {
							style: s,
							classes: ['etFountain'],
						}, false, { time: tMS, space: false }, { x: hD });
					}
				}
			}

			return {
				Still: $Still,
				StraightLine: $StraightLine,
				Rise: $Rise,
				Bounce: $Bounce,
				Speed: $Speed,
				Drop: $Drop,
				Crazy: $Crazy,
				Confetti: $Confetti,
				Throw: $Throw,
				TheCube: $TheCube,
				Fountain: $Fountain,
			};
		}();

		function _queueEmote(e) {
			const sW = window.innerWidth;
			const sH = window.innerHeight;
			const eH = Math.max(cfg.emote.size.min, Math.min(cfg.emote.size.max, Math.floor(sW * cfg.emote.size.ratio.normal), Math.floor(sH * cfg.emote.size.ratio.normal)));
			document.documentElement.style.setProperty('--height', sH + 'px');
			document.documentElement.style.setProperty('--width', sW + 'px');
			const style = cfg.display.styles[shared.random(cfg.display.styles.length)];
			if (style === undefined)
				return;
			display.emote.list[style](e, sW, sH, eH);
		}

		function $showEmotes() {
			if (_tEmote !== false) {
				window.clearTimeout(_tEmote);
				_tEmote = false;
			}
			if (cfg.emote.max > 0 && _eActive >= cfg.emote.max) {
				_tEmote = window.setTimeout(display.emote.showEmotes, 500);
				return;
			}
			let e = null;
			while ((e = _toShow.shift()) !== undefined) {
				_queueEmote(e);
				if (cfg.emote.max > 0 && _eActive > cfg.emote.max) {
					if (cfg.emote.queue > 0 && _toShow.length > cfg.emote.queue)
						_toShow.splice(0, _toShow.length - cfg.emote.queue);
					_tEmote = window.setTimeout(display.emote.showEmotes, 500);
					return;
				}
			}
		}

		function $addToShowList(p) {
			_toShow.push(...p);
		}

		return {
			showEmotes: $showEmotes,
			addToShowList: $addToShowList,
			list: $list,
		};
	}();

	const $kappa = function () {
		const _toKappa = [];
		const _conga = [];
		const _dKappa = 500;

		let _tKappa = false;

		const _list = function () {
			const $Fireworks = function () {
				function $c_Fireworks(kList, sW, sH, eH, iKC) {
					const tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					const oK = kList[shared.random(kList.length)];
					let eW = eH;
					if (oK.hasOwnProperty('width') && oK.hasOwnProperty('height'))
						eW = oK.width / oK.height * eH;
					const eWh = Math.ceil(eW / 2);
					const oX = sW * timing.kappa.Fireworks.origin.x[shared.random(timing.kappa.Fireworks.origin.x.length)] - eWh;
					const oY = sH * timing.kappa.Fireworks.origin.y[shared.random(timing.kappa.Fireworks.origin.y.length)];
					const cX = sW * timing.kappa.Fireworks.dest.x[shared.random(timing.kappa.Fireworks.dest.x.length)];
					const cY = sH * timing.kappa.Fireworks.dest.y[shared.random(timing.kappa.Fireworks.dest.y.length)];
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Fireworks.time);
					const sendUp = Math.floor(tMS * timing.kappa.Fireworks.speed.rocket);
					let s = '--emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eW + 'px;';
					s += ' transform: translate(' + oX + 'px, ' + oY + 'px);';
					const iArr = [];
					_eActive--;
					iArr.push(_addEmoteToDoc(tInit, oK.url, 1, {
						style: s,
						classes: ['ktFireworkRocket'],
					}, true, false, { x: cX - eWh, y: cY }));
					if (oK.hasOwnProperty('zwe')) {
						for (let i = 0, l = oK.zwe.length; i < l; i++) {
							iArr.push(_addEmoteToDoc(tInit, oK.zwe[i].url, 1, {
								style: s,
								classes: ['ktFireworkRocket'],
							}, true, false, { x: cX - eWh, y: cY }));
						}
					}
					window.setTimeout(_explode, sendUp, tInit, kList, iArr, cX, cY, eH, sW, sH, iKC);
				}

				async function _explode(tInit, kList, iArr, cX, cY, eH, sW, sH, iKC) {
					if (_iTitanic > tInit)
						return;
					for (let i = 0, l = iArr.length; i < l; i++) {
						document.body.removeChild(iArr[i]);
						_eActive--;
					}
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Fireworks.time);
					const kTime = Math.floor(tMS * timing.kappa.Fireworks.speed.burst);
					const fA = _kAcTime(iKC, kTime);
					const r = Math.min(sW, sH) * timing.kappa.Fireworks.radius.base;
					const inner = Math.max(3, Math.floor(iKC * timing.kappa.Fireworks.quantity.small));
					const core = Math.floor(iKC * timing.kappa.Fireworks.quantity.medium);
					const outer = Math.max(3, Math.floor(iKC * timing.kappa.Fireworks.quantity.large));
					const lK = kList.length;
					const sR = r * timing.kappa.Fireworks.radius.small;
					for (let v = 0; v < inner; v++) {
						if (_iTitanic > tInit)
							return;
						const sK = kList[shared.random(lK)];
						let eW = eH;
						if (sK.hasOwnProperty('width') && sK.hasOwnProperty('height'))
							eW = sK.width / sK.height * eH;
						const sA = shared.random();
						_eActive--;
						_sparkler(tInit, sK.url, cX, cY, eW, eH, sR, sA);
						if (sK.hasOwnProperty('zwe')) {
							for (let i = 0, l = sK.zwe.length; i < l; i++) {
								_sparkler(tInit, sK.zwe[i].url, cX, cY, eW, eH, sR, sA);
							}
						}
						if (v % fA.ct === fA.ct - 1)
							await _fPause(fA.f);
					}
					await _sleep(Math.floor(tMS * timing.kappa.Fireworks.delays.small));
					const mR = r * timing.kappa.Fireworks.radius.medium;
					const dT = Math.ceil(fA.ct / timing.kappa.Fireworks.spread);
					for (let v = 0; v < core; v++) {
						if (_iTitanic > tInit)
							return;
						const sK = kList[shared.random(lK)];
						let eW = eH;
						if (sK.hasOwnProperty('width') && sK.hasOwnProperty('height'))
							eW = sK.width / sK.height * eH;
						const sA = shared.random();
						_eActive--;
						_sparkler(tInit, sK.url, cX, cY, eW, eH, mR, sA);
						if (sK.hasOwnProperty('zwe')) {
							for (let i = 0, l = sK.zwe.length; i < l; i++) {
								_sparkler(tInit, sK.zwe[i].url, cX, cY, eW, eH, mR, sA);
							}
						}
						if (v % dT === dT - 1)
							await _fPause();
					}
					await _sleep(Math.floor(tMS * timing.kappa.Fireworks.delays.large));
					const lR = r * timing.kappa.Fireworks.radius.large;
					for (let v = 0; v < outer; v++) {
						if (_iTitanic > tInit)
							return;
						const sK = kList[shared.random(lK)];
						let eW = eH;
						if (sK.hasOwnProperty('width') && sK.hasOwnProperty('height'))
							eW = sK.width / sK.height * eH;
						const sA = shared.random();
						_eActive--;
						_sparkler(tInit, sK.url, cX, cY, eW, eH, lR, sA);
						if (sK.hasOwnProperty('zwe')) {
							for (let i = 0, l = sK.zwe.length; i < l; i++) {
								_sparkler(tInit, sK.zwe[i].url, cX, cY, eW, eH, lR, sA);
							}
						}
						if (v % fA.ct === fA.ct - 1)
							await _fPause(fA.f);
					}
				}

				function _sparkler(tInit, url, cX, cY, eW, eH, r, a) {
					if (_iTitanic > tInit)
						return;
					const th = a * _cRadius;
					const eWh = Math.ceil(eW / 2);
					const eX = cX - eWh;
					const hD = Math.floor(eX + r * Math.cos(th));
					const vD = Math.floor(cY + r * Math.sin(th));
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Fireworks.time);
					let s = '--emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eW + 'px;';
					s += ' transform: translate(' + eX + 'px, ' + cY + 'px);';
					s += _styleEmote([], [], [], [], [], [], true, false, true, false, tMS);
					_addEmoteToDoc(tInit, url, 1, {
						style: s,
						classes: ['ktFireworkSparkler'],
					}, false, { space: false, time: tMS }, { x: hD, y: vD });
				}

				return $c_Fireworks;
			}();

			const $Spiral = function () {
				function $c_Spiral(kList, sW, sH, eH, iKC) {
					const tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					const oX = shared.random() * sW;
					const oY = shared.random(sH - eH);
					const r = Math.min(sW, sH);
					shared.doNextFrame(_init, tInit, kList, oX, oY, eH, r, iKC);
				}

				async function _init(tInit, kList, oX, oY, eH, r, iKC) {
					if (_iTitanic > tInit)
						return;
					const c = _cRadius / (_rndFromRange(timing.kappa.Spiral.vectors) + (shared.random() * 2));
					let th = shared.random() * _cRadius;
					const o = shared.random(2) === 0;
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Spiral.time);
					const sA = _kAcTime(iKC, tMS);
					if (sA.ct > timing.kappa.Spiral.bulk)
						sA.ct = timing.kappa.Spiral.bulk;
					for (let i = 0; i < iKC; i++) {
						if (_iTitanic > tInit)
							return;
						if (o) {
							th -= c;
							if (th <= 0)
								th += _cRadius;
						} else {
							th += c;
							if (th >= _cRadius)
								th -= _cRadius;
						}
						const oK = kList[shared.random(kList.length)];
						let eW = eH;
						if (oK.hasOwnProperty('width') && oK.hasOwnProperty('height'))
							eW = oK.width / oK.height * eH;
						const eWh = Math.ceil(eW / 2);
						_eActive--;
						_sparkler(tInit, oK.url, oX - eWh, oY, eW, eH, r, th);
						if (oK.hasOwnProperty('zwe')) {
							for (let j = 0, m = oK.zwe.length; j < m; j++) {
								_sparkler(tInit, oK.zwe[j].url, oX - eWh, oY, eW, eH, r, th);
							}
						}
						if (i % sA.ct === sA.ct - 1)
							await _fPause(sA.f);
					}
				}

				function _sparkler(tInit, url, oX, oY, eW, eH, r, th) {
					if (_iTitanic > tInit)
						return;
					const hD = Math.floor(oX + r * Math.cos(th));
					const vD = Math.floor(oY + r * Math.sin(th));
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Spiral.time);
					let s = '--emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eW + 'px;';
					s += ' transform: translate(' + oX + 'px, ' + oY + 'px);';
					s += _styleEmote([], [], [], [], [], [], true, false, true, false, tMS);
					_addEmoteToDoc(tInit, url, 1, {
						style: s,
						classes: ['ktSpiral'],
					}, false, { space: false, time: tMS }, { x: hD, y: vD });
				}

				return $c_Spiral;
			}();

			const $Pyramid = function () {
				function $c_Pyramid(kList, sW, sH) {
					const tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					const drawn = [];
					let ct = 0;
					const lP = pyramidDist.length;
					const eH = sW / lP;
					for (let i = 0; i < lP; i++) {
						drawn.push(0);
						ct += pyramidDist[i];
					}
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Pyramid.time);
					const sT = tMS * timing.kappa.Pyramid.show.total;
					const tPerB = Math.max(Math.floor(sT / ct), timing.kappa.Pyramid.show.min);
					const eT = tPerB * ct;
					const dT = Math.floor(tMS * timing.kappa.Pyramid.pause);
					let t = 0;
					for (let i = 0; i < ct; i++) {
						if (_iTitanic > tInit)
							return;
						let x;
						do {
							x = shared.random(lP);
						} while (drawn[x] >= pyramidDist[x]);
						const oK = kList[shared.random(kList.length)];
						_block(tInit, oK.url, x, t, eH, sH, drawn[x] + 1, eT + dT);
						if (oK.hasOwnProperty('zwe')) {
							for (let j = 0, l = oK.zwe.length; j < l; j++) {
								_eActive++;
								_block(tInit, oK.zwe[j].url, x, t, eH, sH, drawn[x] + 1, eT + dT);
							}
						}
						drawn[x]++;
						t += tPerB;
					}
				}

				function _block(tInit, url, x, t, eH, sH, dX, aT) {
					if (_iTitanic > tInit)
						return;
					const img = document.createElement('img');
					img.setAttribute('class', 'emote fit ktPyramid');
					_setImgSrc(img, url);
					const h = Math.floor(eH * x);
					const v = -1 * eH;
					const vD = sH - eH * dX;
					let s = 'top: 0px;';
					s += ' left: ' + h + 'px;';
					s += ' --emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eH + 'px;';
					s += ' transform: translateY(' + v + 'px);';
					img.setAttribute('style', s);
					document.body.appendChild(img);
					window.setTimeout(_tDrop, Math.floor(t / 10 + aT), tInit, img, sH);
					window.setTimeout(_tMove, t, tInit, img, vD);
				}

				function _tMove(tInit, img, vD) {
					if (_iTitanic > tInit)
						return;
					img.style.transform = 'translateY(' + vD + 'px)';
				}

				function _tDrop(tInit, img, sH) {
					if (_iTitanic > tInit)
						return;
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Pyramid.time);
					const pT = Math.floor(tMS * timing.kappa.Pyramid.hide);
					img.classList.replace('ktPyramid', 'ktPyramidDrop');
					img.style.transform = 'translateY(' + sH + 'px)';
					_gc.hook(img, false, true, pT);
				}

				return $c_Pyramid;
			}();

			const $SmallPyramid = function () {
				function $c_SmallPyramid(kList, sW, sH) {
					const tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					const drawn = [];
					let ct = 0;
					const lP = pyramidDist.length;
					const eH = Math.min(sW / lP, Math.floor(sW * cfg.emote.size.ratio.small), Math.floor(sH * cfg.emote.size.ratio.small));
					for (let i = 0; i < lP; i++) {
						drawn.push(0);
						ct += pyramidDist[i];
					}
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.SmallPyramid.time);
					const sT = tMS * timing.kappa.SmallPyramid.show.total;
					const tPerB = Math.max(Math.floor(sT / ct), timing.kappa.SmallPyramid.show.min);
					const eT = tPerB * ct;
					const dT = Math.floor(tMS * timing.kappa.SmallPyramid.pause);
					const oX = shared.random(sW - eH * lP);
					let t = 0;
					for (let i = 0; i < ct; i++) {
						if (_iTitanic > tInit)
							return;
						let x;
						do {
							x = shared.random(lP);
						} while (drawn[x] >= pyramidDist[x]);
						const oK = kList[shared.random(kList.length)];
						_block(tInit, oK.url, oX, x, t, eH, sH, drawn[x] + 1, eT + dT);
						if (oK.hasOwnProperty('zwe')) {
							for (let j = 0, l = oK.zwe.length; j < l; j++) {
								_eActive++;
								_block(tInit, oK.zwe[j].url, oX, x, t, eH, sH, drawn[x] + 1, eT + dT);
							}
						}
						drawn[x]++;
						t += tPerB;
					}
				}

				function _block(tInit, url, oX, x, t, eH, sH, dX, aT) {
					if (_iTitanic > tInit)
						return;
					const img = document.createElement('img');
					img.setAttribute('class', 'emote fit ktSmallPyramid');
					_setImgSrc(img, url);
					const h = oX + eH * x;
					const v = -1 * eH;
					const vD = sH - eH * dX;
					let s = 'top: 0px;';
					s += ' left: ' + h + 'px;';
					s += ' --emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eH + 'px;';
					s += ' transform: translateY(' + v + 'px);';
					img.setAttribute('style', s);
					document.body.appendChild(img);
					window.setTimeout(_tDrop, Math.floor(t / 10 + aT), tInit, img, sH);
					window.setTimeout(_tMove, t, tInit, img, vD);
				}

				function _tMove(tInit, img, vD) {
					if (_iTitanic > tInit)
						return;
					img.style.transform = 'translateY(' + vD + 'px)';
				}

				function _tDrop(tInit, img, sH) {
					if (_iTitanic > tInit)
						return;
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.SmallPyramid.time);
					const pT = Math.floor(tMS * timing.kappa.SmallPyramid.hide);
					img.classList.replace('ktSmallPyramid', 'ktSmallPyramidDrop');
					img.style.transform = 'translateY(' + sH + 'px)';
					_gc.hook(img, false, true, pT);
				}

				return $c_SmallPyramid;
			}();

			const $Stampede = function () {
				async function $c_Stampede(kList, sW, sH, eH, iKC) {
					const tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					const bandHeight = eH * timing.kappa.Stampede.height;
					const d = shared.random(2) === 0;
					const bandTop = shared.random(sH - bandHeight + (eH * timing.kappa.Stampede.top.min) + (eH * timing.kappa.Stampede.top.max)) + (eH * (-1 * timing.kappa.Stampede.top.min));
					const b1 = _rndFromRange(timing.kappa.Stampede.bunch[1]);
					const b2 = shared.random(timing.kappa.Stampede.bunch[2] - b1) + b1;
					const b4 = _rndFromRange(timing.kappa.Stampede.bunch[4]);
					_eActive += b1 + b2 + iKC + b4;
					const hasB4 = b4 > 0;
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Stampede.time);
					const tSpeed = Math.floor(tMS * timing.kappa.Stampede.speed);
					const t1 = Math.floor(tSpeed * timing.kappa.Stampede.pause['1']);
					let maxW = 0;
					for (let i = 0, l = kList.length; i < l; i++) {
						let eW = eH;
						if (kList[i].hasOwnProperty('width') && kList[i].hasOwnProperty('height'))
							eW = kList[i].width / kList[i].height * eH;
						if (eW > maxW)
							maxW = eW;
					}
					await _stampede(tInit, kList, b1, t1, false, bandTop, bandHeight, d, sW, eH, maxW);
					if (_iTitanic > tInit)
						return;
					const t2 = Math.floor(tSpeed * timing.kappa.Stampede.pause['2']);
					await _stampede(tInit, kList, b2, t2, false, bandTop, bandHeight, d, sW, eH, maxW);
					if (_iTitanic > tInit)
						return;
					const sA = _kAcTime(iKC, tMS);
					if (sA.ct > timing.kappa.Stampede.maxdensity)
						sA.ct = timing.kappa.Stampede.maxdensity;
					await _stampede(tInit, kList, iKC, hasB4, sA, bandTop, bandHeight, d, sW, eH, maxW);
					if (_iTitanic > tInit)
						return;
					if (hasB4)
						await _stampede(tInit, kList, b4, false, false, bandTop, bandHeight, d, sW, eH, maxW);
				}

				async function _stampede(tInit, kList, ct, pause, sA, bandTop, bandHeight, d, sW, eH, maxW) {
					if (_iTitanic > tInit)
						return;
					const imgs = [];
					for (let i = 0; i < ct; i++) {
						if (_iTitanic > tInit)
							return;
						const oK = kList[shared.random(kList.length)];
						let eW = eH;
						if (oK.hasOwnProperty('width') && oK.hasOwnProperty('height'))
							eW = oK.width / oK.height * eH;
						const y = bandTop + shared.random(bandHeight);
						_eActive--;
						imgs.push(_run(tInit, oK.url, y, d, sW, eW, eH, maxW));
						if (oK.hasOwnProperty('zwe')) {
							for (let j = 0, l = oK.zwe.length; j < l; j++) {
								_run(tInit, oK.zwe[j].url, y, d, sW, eW, eH, maxW);
							}
						}
						if (sA === false)
							await _sleep(_rndFromRange(timing.kappa.Stampede.smallSleep));
						else {
							if (i % sA.ct === sA.ct - 1) {
								let wF = sA.f;
								if (wF === 1)
									wF = shared.random(3);
								else
									wF *= (shared.random() * 3) / 2;
								if (wF !== 0)
									await _fPause(wF);
							}
						}
					}
					if (pause === false)
						return;
					if (pause !== true) {
						await _sleep(pause);
						return;
					}
					do {
						if (_iTitanic > tInit)
							return;
						await _sleep(100);
						for (let i = imgs.length - 1; i >= 0; i--) {
							if (imgs[i] === null || imgs[i].hasAttribute('deleted'))
								imgs.splice(i, 1);
						}
					} while (imgs.length > 0);
				}

				function _run(tInit, url, v, d, sW, eW, eH, maxW) {
					if (_iTitanic > tInit)
						return;
					const h = -2 * maxW;
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Stampede.time);
					const tSpeed = Math.floor(tMS * timing.kappa.Stampede.speed);
					let s = '--emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eW + 'px;';
					if (d)
						s += ' transform: translate(' + sW + 'px, ' + v + 'px);';
					else
						s += ' transform: translate(' + h + 'px, ' + v + 'px);';
					s += _styleEmoteString([], [], [], [], [], []);
					let img;
					if (d)
						img = _addEmoteToDoc(tInit, url, 1, {
							style: s,
							classes: ['ktStampede'],
						}, true, { space: false, time: tSpeed }, { x: h, y: v });
					else
						img = _addEmoteToDoc(tInit, url, 1, {
							style: s,
							classes: ['ktStampede'],
						}, true, { space: false, time: tSpeed }, { x: sW, y: v });
					window.setTimeout(_tMark, tSpeed, tInit, img);
					return img;
				}

				function _tMark(tInit, img) {
					if (_iTitanic > tInit)
						return;
					if (img === null)
						return;
					if (img.parentNode !== null)
						document.body.removeChild(img);
					img.setAttribute('deleted', true);
				}

				return $c_Stampede;
			}();

			const $Conga = function () {
				async function $c_Conga(kList, sW, sH, eH, nM) {
					const tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					let v = 0;
					let unique = false;
					const bS = Math.ceil(eH * timing.kappa.Conga.size);
					const seg = Math.floor(eH * timing.kappa.Conga.height);
					const sht = Math.floor(sH / seg);
					let lns = sht;
					if (nM)
						lns = timing.kappa.Conga.avoidMiddle;
					while (_conga.length >= lns) {
						if (_iTitanic > tInit)
							return;
						await _sleep(250);
					}
					while (!unique) {
						v = shared.random(sht) * seg;
						if (nM) {
							v = shared.random(timing.kappa.Conga.avoidMiddle);
							const hMid = Math.floor(timing.kappa.Conga.avoidMiddle / 2);
							if (v >= hMid)
								v = sht - 1 - (v - hMid);
							v *= seg;
						}
						let found = false;
						for (let i = 0, l = _conga.length; i < l; i++) {
							if (_conga[i].row === v) {
								found = true;
								break;
							}
						}
						if (!found)
							unique = true;
					}
					_conga.push({ row: v, done: false });
					const urls = [];
					const zurls = [];
					const ct = Math.floor(sW / bS);
					for (let i = 0; i < ct; i++) {
						const oK = kList[shared.random(kList.length)];
						urls.push(oK.url);
						const oZ = [];
						if (oK.hasOwnProperty('zwe')) {
							for (let j = 0, l = oK.zwe.length; j < l; j++) {
								oZ.push(oK.zwe[j].url);
							}
						}
						zurls.push(oZ);
					}
					const d = v / seg % 2 === 0;
					const xtra = Math.floor((sW - ct * bS) / 2);
					const imgs = [];
					const zimgs = [];
					for (let i = 0; i < ct; i++) {
						imgs.push(_dance(tInit, urls[i], i, sW, v, eH, bS, ct, d, xtra));
						const oZ = [];
						for (let j = 0, l = zurls[i].length; j < l; j++) {
							oZ.push(_dance(tInit, zurls[i][j], i, sW, v, eH, bS, ct, d, xtra));
						}
						zimgs.push(oZ);
					}
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Conga.time.show);
					await _sleep(tMS);
					let full = false;
					if (_conga.length === sht)
						full = true;
					await _sleep(Math.floor(cfg.display.kappa.conga.time * 1000));
					if (cfg.display.kappa.conga.contagious) {
						let ex = false;
						let lC = _conga.length;
						if (lC > 1)
							ex = true;
						for (let i = 0; i < lC; i++) {
							if (_conga[i].row !== v)
								continue;
							_conga[i].done = true;
							break;
						}
						let done = false;
						while (!done) {
							if (_iTitanic > tInit)
								return;
							lC = _conga.length;
							if (!ex && lC > 1)
								ex = true;
							let notDone = false;
							for (let i = 0; i < lC; i++) {
								if (_conga[i].done === false) {
									notDone = true;
									break;
								}
							}
							if (notDone === false)
								done = true;
							await _sleep(100);
						}
					}
					for (let i = 0, l = imgs.length; i < l; i++) {
						_endDance(tInit, imgs[i], i, sW, v, eH, bS, ct, d, xtra);
						for (let j = 0, m = zimgs[i].length; j < m; j++) {
							_endDance(tInit, zimgs[i][j], i, sW, v, eH, bS, ct, d, xtra);
						}
					}
					await _sleep(tMS);
					for (let i = 0, l = _conga.length; i < l; i++) {
						if (_conga[i].row !== v)
							continue;
						_conga.splice(i, 1);
						break;
					}
				}

				function _dance(tInit, url, col, sW, v, eH, bS, ct, d, xtra) {
					if (_iTitanic > tInit)
						return;
					const box = document.createElement('div');
					box.setAttribute('class', 'scene ktCongaIn');
					const img = document.createElement('img');
					img.setAttribute('class', 'dancer fit');
					_setImgSrc(img, url);
					let s = '--emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eH + 'px;';
					img.setAttribute('style', s);
					let sE = bS * col + xtra;
					let sB = sE - sW;
					if (d) {
						sE = bS * (ct - 1 - col) + xtra;
						sB = sE + sW;
					}
					s = 'top: ' + v + 'px;';
					s += ' left: 0px;';
					s += ' height: ' + bS + 'px;';
					s += ' width: ' + bS + 'px;';
					s += ' z-index: ' + v + ';';
					s += ' transform: translateX(' + sB + 'px);';
					box.setAttribute('style', s);
					_eActive++;
					box.appendChild(img);
					document.body.appendChild(box);
					shared.doNextFrame(_tMove, tInit, box, sE);
					return box;
				}

				function _endDance(tInit, box, col, sW, v, eH, bS, ct, d, xtra) {
					if (_iTitanic > tInit)
						return;
					let sB = bS * col + xtra;
					let sE = sB + sW;
					if (d) {
						sB = bS * (ct - 1 - col) + xtra;
						sE = sB - sW;
					}
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.Conga.time.hide);
					box.classList.replace('ktCongaIn', 'ktCongaOut');
					_gc.hook(box, true, true, Math.ceil(tMS * 1.25));
					shared.doNextFrame(_tMove, tInit, box, sE);
				}

				function _tMove(tInit, box, sE) {
					if (_iTitanic > tInit)
						return;
					box.style.transform = 'translateX(' + sE + 'px)';
				}

				return $c_Conga;
			}();

			const $TheCube = function () {
				function $c_TheCube(kList, sW, sH, eH, bC, iR) {
					const tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					const eHh = Math.ceil(eH / 2);
					const nHh = eHh * -1;
					const sWm = Math.ceil(sW / 2);
					const sHm = Math.ceil(sH / 2);
					const scene = document.createElement('div');
					scene.setAttribute('class', 'scene fit cube kappa');
					const tMS = Math.floor(cfg.emote.time * 1000 * timing.kappa.TheCube.time);
					const cube = document.createElement('div');
					cube.setAttribute('class', 'cube');
					cube.setAttribute('style', 'transform: translateZ(' + nHh + 'px);');
					const sFaces = ['front', 'back', 'right', 'left', 'top', 'bottom'];
					const eFaces = [];
					for (let i = 0; i < 6; i++) {
						eFaces.push(kList[shared.random(kList.length)]);
					}
					for (let i = 0; i < 6; i++) {
						if (!eFaces[i].hasOwnProperty('zwe') || eFaces[i].zwe.length === 0) {
							const iFace = document.createElement('img');
							_setImgSrc(iFace, eFaces[i].url);
							iFace.dataset.face = sFaces[i];
							cube.appendChild(iFace);
						} else {
							const dFace = document.createElement('div');
							dFace.dataset.face = sFaces[i];
							const pFace = document.createElement('img');
							_setImgSrc(pFace, eFaces[i].url);
							dFace.appendChild(pFace);
							for (let j = 0, l = eFaces[i].zwe.length; j < l; j++) {
								const pctZ = document.createElement('img');
								_setImgSrc(pctZ, eFaces[i].zwe[j].url);
								dFace.appendChild(pctZ);
							}
							cube.appendChild(dFace);
						}
					}
					scene.appendChild(cube);
					let h = shared.random(sW - eH);
					let v = shared.random(sH - eH);
					if (bC) {
						h = Math.floor(sWm - eHh);
						v = Math.floor(sHm - eHh);
					}
					let s = '--emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eH + 'px;';
					s += ' --cube-depth: ' + eHh + 'px;';
					s += ' perspective: ' + eH * 3 + 'px;';
					s += ' transform: translate(' + h + 'px, ' + v + 'px);';
					s += _styleEmote([], [], [], [], [], [], cfg.emote.in.fade, false, cfg.emote.out.fade, false, tMS);
					scene.setAttribute('style', s);
					document.body.appendChild(scene);
					_gc.hook(scene, false, 6, tMS);
					shared.doNextFrame(_tMove, tInit, cube, iR, eH);
				}

				function _tMove(tInit, cube, iR, eH) {
					if (_iTitanic > tInit)
						return;
					const nHh = Math.ceil(eH / 2) * -1;
					let rX = 0;
					let rY = 0;
					while (Math.abs(rX) + Math.abs(rY) < 45) {
						rX = (360 - shared.random() * 720) * iR;
						rY = (360 - shared.random() * 720) * iR;
					}
					cube.style.transform = 'translateZ(' + nHh + 'px) rotateX(' + rX + 'deg) rotateY(' + rY + 'deg)';
				}

				return $c_TheCube;
			}();

			const $Text = function () {
				let _mL = 0;

				function $c_Text(kList, sW, sH, sMsg, iTime) {
					const tInit = new Date().getTime();
					if (_iTitanic > tInit)
						return;
					const msgDist = _buildMsgArr(sMsg);
					let ct = 0;
					let ctT = 0;
					const drawn = [];
					const lM = msgDist.length;
					for (let x = 0; x < lM; x++) {
						const lX = msgDist[x].length;
						for (let y = 0; y < lX; y++) {
							if (msgDist[x][y] !== 0)
								ctT += 1;
						}
						ct += lX;
						drawn.push(0);
					}
					const eH = Math.min(Math.floor(sW / (lM + 2)), Math.floor(sW * cfg.emote.size.ratio.small), Math.floor(sH * cfg.emote.size.ratio.small));
					const tMS = Math.floor(iTime * 1000 * timing.kappa.Text.time);
					const sT = tMS * timing.kappa.Text.show.total;
					const tPerB = Math.max(Math.floor(sT / ct), timing.kappa.Text.show.min);
					const eT = tPerB * ctT;
					const lF = msgDist[0].length;
					const lFS = eH * lF;
					const vH = shared.random(sH - lFS) + lFS;
					const oX = shared.random(sW - eH * lM);
					let t = 0;
					for (let i = 0; i < ct; i++) {
						if (_iTitanic > tInit)
							return;
						let x;
						do {
							x = shared.random(lM);
						} while (drawn[x] >= msgDist[x].length);
						if (msgDist[x][drawn[x]] !== 0) {
							const oK = kList[shared.random(kList.length)];
							_block(tInit, oK.url, vH, oX, drawn[x] + 1, tPerB, eT, iTime, x, t, sH, eH);
							if (oK.hasOwnProperty('zwe')) {
								for (let j = 0, l = oK.zwe.length; j < l; j++) {
									_block(tInit, oK.zwe[j].url, vH, oX, drawn[x] + 1, tPerB, eT, iTime, x, t, sH, eH);
								}
							}
							t += tPerB;
						}
						drawn[x]++;
					}
				}

				function _buildMsgArr(s) {
					const o = [];
					const spc = [];
					if (_mL === 0) {
						for (let i = 0, k = Object.keys(alnumDist), l = k.length; i < l; i++) {
							_mL = Math.max(_mL, alnumDist[k[i]][0].length);
						}
					}
					for (let y = 0; y < _mL; y++) {
						spc.push(0);
					}
					for (let i = 0, l = s.length; i < l; i++) {
						if (i > 0)
							o.push(spc);
						if (s[i] === ' ') {
							o.push(spc);
							o.push(spc);
							continue;
						}
						const v = s[i];
						if (!alnumDist.hasOwnProperty(v))
							continue;
						const c = alnumDist[v];
						for (let x = 0, m = c.length; x < m; x++) {
							o.push(c[x]);
						}
					}
					return o;
				}

				function _block(tInit, url, vH, oX, dX, tPerB, eT, iTime, x, t, sH, eH) {
					if (_iTitanic > tInit)
						return;
					const img = document.createElement('img');
					img.setAttribute('class', 'emote fit');
					_setImgSrc(img, url);
					const h = oX + eH * x;
					const v = -1 * eH;
					const vD = vH - eH * dX;
					let s = 'top: 0px;';
					s += ' left: ' + h + 'px;';
					s += ' --emote-height: ' + eH + 'px;';
					s += ' --emote-width: ' + eH + 'px;';
					s += ' transition: transform ' + tPerB + 'ms ease-in;';
					s += ' transform: translateY(' + v + 'px);';
					img.setAttribute('style', s);
					document.body.appendChild(img);
					_eActive++;
					const tMS = Math.floor(iTime * 1000 * timing.kappa.Text.time);
					window.setTimeout(_tDrop, Math.floor(eT + tMS + t / 10), tInit, img, sH, tMS);
					window.setTimeout(_tMove, t, tInit, img, vD);
				}

				function _tMove(tInit, img, vD) {
					if (_iTitanic > tInit)
						return;
					img.style.transform = 'translateY(' + vD + 'px)';
				}

				function _tDrop(tInit, img, sH, tMS) {
					if (_iTitanic > tInit)
						return;
					const pT = Math.floor(tMS * timing.kappa.Text.hide);
					img.style.transform = 'translateY(' + sH + 'px)';
					img.style.transitionDuration = pT + 'ms';
					_gc.hook(img, false, true, pT);
				}

				return $c_Text;
			}();

			return {
				Fireworks: $Fireworks,
				Spiral: $Spiral,
				Pyramid: $Pyramid,
				SmallPyramid: $SmallPyramid,
				Stampede: $Stampede,
				Conga: $Conga,
				TheCube: $TheCube,
				Text: $Text,
			};
		}();

		function _canShowKappa(k) {
			if (cfg.emote.max < 1)
				return true;
			if (_eActive < 1)
				return true;
			let tC = cfg.display.kappa.count;
			if (k !== false)
				tC = _getKappaCountEstimate(k);
			const cM = Math.max(cfg.emote.max, tC);
			return _eActive + tC < cM;
		}

		function _getKappaCountParam(p) {
			const a = p.split(' ');
			for (let i = 0, l = a.length; i < l; i++) {
				if (!isNaN(a[i]))
					return parseInt(a[i], 10);
			}
			return false;
		}

		function _getKappaCountEstimate(k) {
			switch (k.style) {
				case 'Pyramid':
				case 'SmallPyramid':
					let c = 0;
					for (let i = 0, l = pyramidDist.length; i < l; i++) {
						c += pyramidDist[i];
					}
					return c;
				case 'Fireworks':
					const inner = Math.max(3, Math.floor(k.count * timing.kappa.Fireworks.quantity.small));
					const core = Math.floor(k.count * timing.kappa.Fireworks.quantity.medium);
					const outer = Math.max(3, Math.floor(k.count * timing.kappa.Fireworks.quantity.large));
					return 1 + inner + core + outer;
				case 'Conga':
					const sW = window.innerWidth;
					const sH = window.innerHeight;
					const eH = Math.max(cfg.emote.size.min, Math.min(cfg.emote.size.max, Math.floor(sW * cfg.emote.size.ratio.normal), Math.floor(sH * cfg.emote.size.ratio.normal)));
					const bS = Math.ceil(eH * timing.kappa.Conga.size);
					return Math.floor(sW / bS);
				case 'TheCube':
					return 6;
			}
			return k.count;
		}

		async function $show(kList, kStyle) {
			const sW = window.innerWidth;
			const sH = window.innerHeight;
			const eH = Math.max(cfg.emote.size.min, Math.min(cfg.emote.size.max, Math.floor(sW * cfg.emote.size.ratio.normal), Math.floor(sH * cfg.emote.size.ratio.normal)));
			const eHh = Math.max(cfg.emote.size.min, Math.min(Math.floor(cfg.emote.size.max / 2), Math.floor(sW * cfg.emote.size.ratio.small), Math.floor(sH * cfg.emote.size.ratio.small)));
			const sB = sH - eH;
			document.documentElement.style.setProperty('--height', sH + 'px');
			document.documentElement.style.setProperty('--width', sW + 'px');
			const waitFor = _getKappaCountEstimate(kStyle);

			if (!_canShowKappa(kStyle)) {
				_toKappa.push({
					list: kList,
					style: kStyle.style,
					prefs: kStyle.prefs,
					params: kParams,
				});
				if (_tKappa !== false) {
					window.clearTimeout(_tKappa);
					_tKappa = false;
				}
				_tKappa = window.setTimeout(_showKappas, _dKappa);
				return;
			}

			_eActive += waitFor;
			const lK = kList.length;
			const tInit = new Date().getTime();
			let estMS = Math.floor(cfg.emote.time * 1000);
			if (timing.kappa.hasOwnProperty(kStyle.style) && timing.kappa[kStyle.style].hasOwnProperty('time'))
				estMS = Math.floor(cfg.emote.time * 1000 * timing.kappa[kStyle.style].time);
			else if (timing.display.hasOwnProperty(kStyle.style) && timing.display[kStyle.style].hasOwnProperty('time'))
				estMS = Math.floor(cfg.emote.time * 1000 * timing.display[kStyle.style].time);

			switch (kStyle.style) {
				case 'Stampede':
					_eActive -= waitFor;
					await _list.Stampede(kList, sW, sH, eH, kStyle.count);
					break;
				case 'Fireworks':
					_list.Fireworks(kList, sW, sH, eHh, kStyle.count);
					break;
				case 'Spiral':
					_list.Spiral(kList, sW, sH, eHh, kStyle.count);
					break;
				case 'Pyramid':
					_list.Pyramid(kList, sW, sH);
					break;
				case 'SmallPyramid':
					_list.SmallPyramid(kList, sW, sH);
					break;
				case 'Conga':
					_eActive -= waitFor;
					let avoidMiddle = false;
					if (cfg.display.kappa.conga.hasOwnProperty('avoidMiddle') && cfg.display.kappa.conga.avoidMiddle === true)
						avoidMiddle = true;
					if (kStyle.prefs.hasOwnProperty('avoidMiddle') && kStyle.prefs.avoidMiddle === true)
						avoidMiddle = true;
					_list.Conga(kList, sW, sH, eH, avoidMiddle);
					break;
				case 'Text':
					_eActive -= waitFor;
					let sTM = 'HYPE!';
					if (cfg.display.kappa.styles.hasOwnProperty(kStyle.style) && cfg.display.kappa.styles[kStyle.style].hasOwnProperty('message'))
						sTM = cfg.display.kappa.styles[kStyle.style].message[shared.random(cfg.display.kappa.styles[kStyle.style].message.length)];
					if (kStyle.prefs.hasOwnProperty('message') && Array.isArray(kStyle.prefs.message) && kStyle.prefs.message.length > 0)
						sTM = kStyle.prefs.message[shared.random(kStyle.prefs.message.length)];
					let sTT = cfg.emote.time;
					if (cfg.display.kappa.styles.hasOwnProperty(kStyle.style) && cfg.display.kappa.styles[kStyle.style].hasOwnProperty('time'))
						sTT = cfg.display.kappa.styles[kStyle.style].time;
					if (kStyle.prefs.hasOwnProperty('time') && kStyle.prefs.time > 0)
						sTT = kStyle.prefs.time;
					_list.Text(kList, sW, sH, sTM, sTT);
					break;
				case 'TheCube':
					const cS = Math.min(sW, sH);
					let sCS = 8 / 10;
					if (cfg.display.kappa.styles.hasOwnProperty(kStyle.style) && cfg.display.kappa.styles[kStyle.style].hasOwnProperty('size'))
						sCS = cfg.display.kappa.styles[kStyle.style].size;
					if (kStyle.prefs.hasOwnProperty('size'))
						sCS = kStyle.prefs.size;
					let sCC = true;
					if (cfg.display.kappa.styles.hasOwnProperty(kStyle.style) && cfg.display.kappa.styles[kStyle.style].hasOwnProperty('center'))
						sCC = cfg.display.kappa.styles[kStyle.style].center;
					if (kStyle.prefs.hasOwnProperty('center'))
						sCC = kStyle.prefs.center;
					let sCR = 5;
					if (cfg.display.kappa.styles.hasOwnProperty(kStyle.style) && cfg.display.kappa.styles[kStyle.style].hasOwnProperty('rotations'))
						sCR = cfg.display.kappa.styles[kStyle.style].rotations;
					if (kStyle.prefs.hasOwnProperty('rotations'))
						sCR = kStyle.prefs.rotations;
					let bF = false;
					if (cfg.display.kappa.styles.hasOwnProperty(kStyle.style) && cfg.display.kappa.styles[kStyle.style].hasOwnProperty('faces'))
						bF = cfg.display.kappa.styles[kStyle.style].faces;
					if (kStyle.prefs.hasOwnProperty('faces'))
						bF = kStyle.prefs.faces === true;
					let kUse = [];
					if (bF)
						kUse = kList;
					else
						kUse.push(kList[shared.random(lK)]);
					_list.TheCube(kUse, sW, sH, Math.floor(cS * sCS), sCC, sCR);
					break;
				case 'Burst':
					const oH = _rndFromRange(timing.kappa[kStyle.style].left);
					const oV = _rndFromRange(timing.kappa[kStyle.style].top) * sB;
					const bA = _kAcTime(kStyle.count, estMS);
					for (let i = 0; i < kStyle.count; i++) {
						if (_iTitanic > tInit)
							return;
						const rB = shared.random(lK);
						_eActive--;
						let eWb = eH;
						if (kList[rB].hasOwnProperty('width') && kList[rB].hasOwnProperty('height'))
							eWb = kList[rB].width / kList[rB].height * eH;
						const sRb = sW - Math.ceil(eWb / 2);
						display.emote.list.StraightLine(kList[rB], sW, sH, eH, oH * sRb, oV, false, tInit);
						if (i % bA.ct === bA.ct - 1)
							await _fPause(bA.f);
					}
					break;
				case 'Fountain':
					const fX = _rndFromRange(timing.kappa[kStyle.style].left) * sW;
					const fY = _rndFromRange(timing.kappa[kStyle.style].top);
					const fA = _kAcTime(kStyle.count, estMS);
					for (let i = 0; i < kStyle.count; i++) {
						if (_iTitanic > tInit)
							return;
						const rF = shared.random(lK);
						_eActive--;
						display.emote.list.Fountain(kList[rF], sW, sH, eH, fX, fY, false, tInit);
						if (i % fA.ct === fA.ct - 1)
							await _fPause(fA.f);
					}
					break;
				case 'Confetti':
					const cA = _kAcTime(kStyle.count, estMS);
					for (let i = 0; i < kStyle.count; i++) {
						if (_iTitanic > tInit)
							return;
						const rN = shared.random(lK);
						_eActive--;
						display.emote.list.Confetti(kList[rN], sW, sH, eHh, false, tInit);
						if (i % cA.ct === cA.ct - 1)
							await _fPause(cA.f);
					}
					break;
			}
		}

		function $hide() {
			if (_tKappa !== false) {
				window.clearTimeout(_tKappa);
				_tKappa = false;
			}
			_toKappa.length = 0;
			_conga.length = 0;
		}

		function _sleep(ms) {
			if (ms < shared.mspf.value)
				return _fPause();
			return new Promise(
				function (resolve) {
					let n = 0;

					function _next(ts) {
						if (n === 0) {
							n = ts;
							window.requestAnimationFrame(_next);
							return;
						} else if (ts - n < ms) {
							window.requestAnimationFrame(_next);
							return;
						}
						resolve(true);
					}

					window.requestAnimationFrame(_next);
				},
			);
		}

		function _fPause(frames = 1) {
			return new Promise(
				function (resolve) {
					if (frames < 1) {
						resolve(false);
						return;
					}
					let n = 0;

					function _next() {
						n++;
						if (n < frames) {
							window.requestAnimationFrame(_next);
							return;
						}
						resolve(true);
					}

					window.requestAnimationFrame(_next);
				},
			);
		}

		function _kAcTime(ct, t = false) {
			if (t === false)
				t = Math.floor(cfg.emote.time * 1000);
			const f = Math.floor(t / shared.mspf.value);
			const r = f / ct;
			if (r > 1)
				return { f: Math.ceil(r), ct: 1 };
			return { f: 1, ct: Math.ceil(1 / r) };
		}

		return {
			show: $show,
			hide: $hide,
		};
	}();

	const _gc = function () {
		const _toGC = {};

		let _tGC = false;

		function _doGC() {
			if (_tGC === false)
				return;
			window.clearTimeout(_tGC);
			_tGC = false;
			let done = true;
			const tNow = new Date().getTime();
			for (const idx in _toGC) {
				if (!_toGC.hasOwnProperty(idx))
					continue;
				done = false;
				const i = _toGC[idx].img;
				const t = _toGC[idx].end;
				const d = _toGC[idx].dec;
				if (_toGC[idx].space) {
					const r = i.getBoundingClientRect();
					if (t > tNow && r.bottom > 0 && r.right > 0 && r.top < window.innerHeight && r.left < window.innerWidth)
						continue;
				} else {
					if (t > tNow)
						continue;
				}
				delete _toGC[idx];
				if (i.parentNode !== null)
					document.body.removeChild(i);
				if (d === true)
					_eActive--;
				else if (d !== false && !isNaN(d))
					_eActive -= d;
			}
			if (!done)
				_tGC = window.setTimeout(_doGC, 500);
		}

		function $hook(img, space = true, decActive = true, t = false) {
			if (t === false)
				t = Math.floor(cfg.emote.time * 1000);
			let x = 0;
			do {
				x++;
			} while (_toGC.hasOwnProperty(x));
			_toGC[x] = { img: img, space: space, dec: decActive, end: new Date().getTime() + t };
			if (_tGC === false)
				_tGC = window.setTimeout(_doGC, 500);
		}

		return {
			hook: $hook,
		};
	}();

	function _rndFromRange(range) {
		return shared.random(range.max - range.min) + range.min;
	}

	function _safePoints(h, v, th, r, fL, fT, fR, fB) {
		const hD = Math.floor(h + r * Math.cos(th));
		const vD = Math.floor(v + r * Math.sin(th));
		const slope = Math.tan(th);
		let hL = Number.MAX_SAFE_INTEGER;
		let hU = 0;
		let vL = Number.MAX_SAFE_INTEGER;
		let vU = 0;
		if (hD < fL)
			hL = h - (hU = fL);
		else if (hD > fR)
			hL = (hU = fR) - h;
		if (vD < fT)
			vL = v - (vU = fT);
		else if (vD > fB)
			vL = (vU = fB) - v;
		if (vU === 0 && hU === 0)
			return true;
		let vT = vU;
		let hT = hU;
		if (hL > vL)
			hT = Math.floor((vT - v) / slope + h);
		else
			vT = Math.floor((hT - h) * slope + v);
		const l = Math.sqrt(Math.abs(h - hT) ** 2 + Math.abs(v - vT) ** 2);
		return (l > Math.ceil(r / 2));
	}

	function _addEmoteToDoc(tInit, uri, variationSize, attrs = {}, r = false, oGC = {}, oT = false) {
		if (_iTitanic > tInit)
			return;
		const img = document.createElement('img');
		const c = [];
		c.push('emote');
		if (attrs.hasOwnProperty('classes'))
			c.push(...attrs.classes);
		const rV = variationSize.toFixed(3).replace('.', '_');
		c.push('eSize-' + rV);
		img.classList.add(...c);
		_setImgSrc(img, uri);
		if (attrs.hasOwnProperty('style'))
			img.setAttribute('style', attrs.style);
		if (attrs.hasOwnProperty('dataset')) {
			for (let i = 0, k = Object.keys(attrs.dataset), l = k.length; i < l; i++) {
				img.setAttribute('data-' + k[i], attrs.dataset[k[i]]);
			}
		}
		_eActive++;
		document.body.appendChild(img);
		let space = true;
		let decActive = true;
		let t = false;
		if (oGC !== false) {
			if (oGC.hasOwnProperty('space'))
				space = oGC.space;
			if (oGC.hasOwnProperty('decrement'))
				decActive = oGC.decrement;
			if (oGC.hasOwnProperty('time'))
				t = oGC.time;
			_gc.hook(img, space, decActive, t);
		}
		if (oT !== false) {
			let sTF = null;
			if (oT.hasOwnProperty('x') && oT.hasOwnProperty('y'))
				sTF = 'translate(' + oT.x + 'px, ' + oT.y + 'px)';
			else if (oT.hasOwnProperty('x'))
				sTF = 'translateX(' + oT.x + 'px)';
			else if (oT.hasOwnProperty('y'))
				sTF = 'translateY(' + oT.y + 'px)';
			if (sTF !== null)
				shared.doNextFrame(_tMoveOnDock, tInit, img, sTF);
		}
		if (r)
			return img;
	}

	function _tMoveOnDock(tInit, img, sTF) {
		if (_iTitanic > tInit)
			return;
		img.style.transform = sTF;
	}

	function _setImgSrc(img, url) {
		img.alt = '';
		img.onload = function () {
			img.onerror = null;
			img.onload = null;
		};
		img.onerror = function () {
			img.onerror = null;
			img.onload = null;
			img.src = bareList[shared.random(bareList.length)].url;
		};
		img.src = url;
	}

	function _styleEmoteString(aNames, aDelays, aDurs, aTimings, aFills, aIters) {
		let s = '';
		if (aNames.length > 0) {
			s += ' animation-name: ' + aNames.join() + ';';
			s += ' animation-delay: ' + aDelays.join() + ';';
			s += ' animation-duration: ' + aDurs.join() + ';';
			s += ' animation-timing-function: ' + aTimings.join() + ';';
			s += ' animation-fill-mode: ' + aFills.join() + ';';
			s += ' animation-iteration-count: ' + aIters.join() + ';';
		}
		if (aNames.includes('fadeIn'))
			s += ' opacity: 0;';
		return s;
	}

	function _styleEmote(aNames, aDelays, aDurs, aTimings, aFills, aIters, fadeIn = true, zoomIn = true, fadeOut = true, zoomOut = true, tMS = false) {
		if (tMS === false)
			tMS = Math.floor(cfg.emote.time * 1000);
		const tFI = _tAnim.fade.in / 100;
		const tFO = _tAnim.fade.out / 100;
		const tZI = _tAnim.zoom.in / 100;
		const tZO = _tAnim.zoom.out / 100;
		if (fadeIn) {
			aNames.push('fadeIn');
			aDelays.push('0s');
			aDurs.push(Math.floor(tMS * tFI) + 'ms');
			aTimings.push('ease-in');
			aFills.push('forwards');
			aIters.push('1');
		}
		if (zoomIn) {
			aNames.push('zoomIn');
			aDelays.push('0s');
			aDurs.push(Math.floor(tMS * tZI) + 'ms');
			aTimings.push('linear');
			aFills.push('forwards');
			aIters.push('1');
		}
		if (fadeOut) {
			aNames.push('fadeOut');
			aDelays.push(Math.floor(tMS - tMS * tFO) + 'ms');
			aDurs.push(Math.floor(tMS * tFO) + 'ms');
			aTimings.push('ease-out');
			aFills.push('forwards');
			aIters.push('1');
		} else {
			aNames.push('noFadeOut');
			aDelays.push(tMS - 50 + 'ms');
			aDurs.push('50ms');
			aTimings.push('ease-out');
			aFills.push('forwards');
			aIters.push('1');
		}
		if (zoomOut) {
			aNames.push('zoomOut');
			aDelays.push(Math.floor(tMS - tMS * tZO) + 'ms');
			aDurs.push(Math.floor(tMS * tZO) + 'ms');
			aTimings.push('linear');
			aFills.push('forwards');
			aIters.push('1');
		}
		return _styleEmoteString(aNames, aDelays, aDurs, aTimings, aFills, aIters);
	}

	function $eraseAll() {
		_iTitanic = new Date().getTime();
		display.kappa.hide();
		const cubes = document.getElementsByClassName('scene');
		while (cubes.length) {
			cubes[0].parentElement.removeChild(cubes[0]);
		}
		const imgs = document.getElementsByTagName('img');
		while (imgs.length) {
			imgs[0].parentElement.removeChild(imgs[0]);
		}
		service.parse.clearCooldowns();
		_eActive = 0;
	}

	return {
		gc: _gc,
		emote: $emote,
		kappa: $kappa,
		eraseAll: $eraseAll,
	};
}();

const shared = function () {
	function $random(m) {
		const r = new Uint32Array(1);
		window.crypto.getRandomValues(r);
		const f = r[0] / 4294967295;
		if (m === undefined)
			return f;
		if (m < 1)
			return f * m;
		return Math.floor(f * m);
	}

	function $doNextFrame(cb) {
		const a = [];
		for (let i = 1, l = arguments.length; i < l; i++) {
			a.push(arguments[i]);
		}
		let n = false;

		function _next() {
			if (n === false) {
				n = true;
				window.requestAnimationFrame(_next);
				return;
			}
			cb(...a);
		}

		window.requestAnimationFrame(_next);
	}

	const $mspf = function () {
		let _init = 0;
		const _avg = [];

		function $init() {
			if (shared.mspf.value !== 0)
				return;
			window.requestAnimationFrame(_test);
		}

		function _test(ms) {
			if (_init !== 0)
				_avg.push(ms - _init);
			if (_avg.length > 2)
				shared.mspf.value = _avg.reduce((a, b) => (a + b)) / _avg.length;
			if (_avg.length > 300) {
				shared.mspf.value = Math.round(shared.mspf.value * 1000) / 1000;
				return;
			}
			_init = ms;
			window.requestAnimationFrame(_test);
		}

		return {
			init: $init,
			value: 0,
		};
	}();

	return {
		random: $random,
		doNextFrame: $doNextFrame,
		mspf: $mspf,
	};
}();

const startup = function () {
	async function $c_startup() {
		shared.mspf.init();
		if (cfg.display.hasOwnProperty('hue') && cfg.display.hue === 'rave' && !_hasRaveToggle())
			document.documentElement.classList.add('rave');
		_css();
	}

	return $c_startup;
}();

window.startup = startup;
window.emote = display.emote;
window.kappagen = display.kappa;
window.random = shared.random;
</script>

<style>
.etStraightLine {
  transition: transform linear 5000ms;
}
.etThrowTwist {
  transition: transform cubic-bezier(0.32, 0, 0.67, 0) 700ms;
}
.etThrowDrop {
  transition: transform cubic-bezier(0.5, 0, 0.75, 0) 4000ms;
}
.etFountain {
  transition: transform cubic-bezier(0, 0, 0.58, 1) 2500ms;
}
div.scene.cube:not(.kappa),
div.scene.cube:not(.kappa) div.cube {
  transition: transform linear 5000ms;
}
.ktFireworkRocket {
  transition: transform linear 2000ms;
}
.ktFireworkSparkler {
  transition: transform ease-out 5000ms;
}
.ktSpiral {
  transition: transform ease-out 2500ms;
}
.ktPyramid {
  transition: transform ease-in 75ms;
}
.ktPyramidDrop {
  transition: transform ease-in 50ms;
}
.ktSmallPyramid {
  transition: transform ease-in 100ms;
}
.ktSmallPyramidDrop {
  transition: transform ease-in 50ms;
}
.ktStampede {
  transition: transform linear 2000ms;
}
.ktCongaIn {
  transition: transform linear 10000ms;
}
.ktCongaOut {
  transition: transform ease-in 10000ms;
}
div.scene.cube.kappa,
div.scene.cube.kappa div.cube {
  transition: transform linear 5000ms;
}

[data-squash="vertical"] {
  transform: scale(2, 0.7);
}
[data-squash="horizontal"] {
  transform: scale(0.7, 2);
}
[data-squash="no"] {
  transform: scale(1, 1);
}

[data-origin="center"] {
  transform-origin: center center;
}

[data-origin="topleft"] {
  transform-origin: left top;
}

[data-origin="topright"] {
  transform-origin: right top;
}

[data-origin="top"] {
  transform-origin: center top;
}

[data-origin="bottom"] {
  transform-origin: center bottom;
}

[data-origin="left"] {
  transform-origin: left center;
}

[data-origin="right"] {
  transform-origin: right center;
}
</style>
