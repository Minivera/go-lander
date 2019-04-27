package main

import lander "github.com/minivera/go-lander"

func wrapper(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("section", map[string]string{}, []lander.Node{
			lander.Html("div", map[string]string{}, children).Style(`
				margin-left: auto;
				margin-right: auto;
				max-width: 600px;
			`),
		}).Style(`
			padding-left: 2rem;
			padding-right: 2rem;
			width: 100%;
		`),
	}
}

func manifoldContainer(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("div", map[string]string{
			"itemscope": "true",
			"itemtype":  "https://schema.org/Brand",
		}, children).Style(`
			display: inline-flex;
			height: 3rem;
			width: auto;
		`),
	}
}

func manifold(_ map[string]string, _ []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Component(manifoldContainer, map[string]string{}, []lander.Node{
			lander.Html("div", map[string]string{}, []lander.Node{
				lander.Svg("svg", map[string]string{
					"height":      "3rem",
					"width":       "3rem",
					"viewBox":     "0 0 512 512",
					"xmlns":       "http://www.w3.org/2000/svg",
					"xmlns:xlink": "http://www.w3.org/1999/xlink",
				}, []lander.Node{
					lander.Svg("title", map[string]string{
						"itemProp": "name",
					}, []lander.Node{
						lander.Text("Manifold"),
					}),
					lander.Svg("defs", map[string]string{}, []lander.Node{
						lander.Svg("linearGradient#manifold-a", map[string]string{
							"x1": "1.202%",
							"y1": "95.992%",
							"x2": "92.841%",
							"y2": "42.213%",
						}, []lander.Node{
							lander.Svg("stop", map[string]string{
								"stop-color": "#FF0264",
								"offset":     "0%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#FE624E",
								"offset":     "40.83%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#FDBC39",
								"offset":     "81.65%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#FDDF31",
								"offset":     "100%",
							}, []lander.Node{}),
						}),
						lander.Svg("linearGradient#manifold-b", map[string]string{
							"x1": ".042%",
							"y1": "49.985%",
							"x2": "100.019%",
							"y2": "49.985%",
						}, []lander.Node{
							lander.Svg("stop", map[string]string{
								"stop-color": "#140A3B",
								"offset":     "0%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#064E91",
								"offset":     "68.04%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#006AB4",
								"offset":     "100%",
							}, []lander.Node{}),
						}),
						lander.Svg("linearGradient#manifold-c", map[string]string{
							"x1": ".465%",
							"y1": "50.005%",
							"x2": "173.382%",
							"y2": "50.005%",
						}, []lander.Node{
							lander.Svg("stop", map[string]string{
								"stop-color": "#349FD3",
								"offset":     "0%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#218ABE",
								"offset":     "35.17%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#00679B",
								"offset":     "100%",
							}, []lander.Node{}),
						}),
						lander.Svg("linearGradient#manifold-d", map[string]string{
							"x1": "-.029%",
							"y1": "49.982%",
							"x2": "99.865%",
							"y2": "49.982%",
						}, []lander.Node{
							lander.Svg("stop", map[string]string{
								"stop-color": "#A34CB4",
								"offset":     "0%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#954DB2",
								"offset":     "9.192%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#3C52A1",
								"offset":     "71.69%",
							}, []lander.Node{}),
							lander.Svg("stop", map[string]string{
								"stop-color": "#19549B",
								"offset":     "100%",
							}, []lander.Node{}),
						}),
					}),
					lander.Svg("g", map[string]string{
						"fill":      "none",
						"fill-rule": "evenodd",
					}, []lander.Node{
						lander.Svg("mask#manifold-mask", map[string]string{
							"fill": "#fff",
						}, []lander.Node{
							lander.Svg("circle", map[string]string{
								"cx": "256",
								"cy": "256",
								"r":  "256",
							}, []lander.Node{}),
						}),
						lander.Svg("path", map[string]string{
							"fill": "url(#manifold-a)",
							"mask": "url(#manifold-mask)",
							"d":    "M512 305.672H0V0h512z",
						}, []lander.Node{}),
						lander.Svg("path", map[string]string{
							"fill":      "url(#manifold-b)",
							"fill-rule": "nonzero",
							"mask":      "url(#manifold-mask)",
							"d":         "M512 437.005C405.284 568.17 252.826 456.57 252.826 456.57L85.474 307.623l39.776-34.909c19.464-17.349 17.56-138.578 75.742-161.216 25.177-9.733 51.834 2.538 56.912 4.23 6.064 2.022 165.148 63.245 254.096 97.623v223.654z",
						}, []lander.Node{}),
						lander.Svg("path", map[string]string{
							"fill":      "url(#manifold-c)",
							"fill-rule": "nonzero",
							"mask":      "url(#manifold-mask)",
							"d":         "M512 214.777V512H212.82c-6.33-25.985-12.886-62.758-12.886-86.955 0-71.511 37.448-119.114 76.377-98.38 10.578 5.712 19.887 17.56 30.254 33.64 15.868 24.964 32.37 37.024 47.815 37.447 5.29.212 10.79-1.058 15.656-3.597 11.214-5.924 21.369-18.618 29.409-38.082 5.924-13.752 10.578-31.101 13.963-51.412l1.693-10.79c4.866-30.043 12.906-53.104 22.85-68.972 0 0 .423-.635 1.057-1.48 4.655-6.771 23.908-29.832 67.914-11.002 1.653.772 3.35 1.56 5.078 2.36z",
						}, []lander.Node{}),
						lander.Svg("path", map[string]string{
							"fill":      "url(#manifold-d)",
							"fill-rule": "nonzero",
							"mask":      "url(#manifold-mask)",
							"d":         "M0 213.69l269.329 109.378C233.997 311.43 204.8 367.92 204.8 429.91c0 24.965 5.712 53.104 12.483 77.646 0 .437.305 1.918.917 4.443H0V213.69z",
						}, []lander.Node{}),
					}),
				}),
			}).Style(`
				align-items: center;
				display: flex;
				flex: 512;
				height: 100%;
			`).SelectorStyle(" svg", `
				fill: currentColor;
				height: 100%;
				width: auto;
			`),
		}),
	}
}

const (
	baseButtonStyle = `
		align-items: center;
		background: white;
		border-radius: 4px;
		border: 1px solid #C4C7CA;
		box-shadow: 0 2px 6px rgba(31, 31, 38, .1);
		color: #323940;
		cursor: pointer;
		display: inline-flex;
		font-size: 15px;
		font-weight: 500;
		min-height: 2.5rem;
		justify-content: center;
		line-height: 1;
		padding-bottom: 0.75em;
		padding-left: 1.5em;
		padding-right: 1.5em;
		padding-top: 0.75em;
		text-align: center;
		text-decoration: none;
		transition: background-color 200ms, border-color 150ms, color 200ms;
		white-space: nowrap;
	`
	focusHoverButtonStyle = `
		border-color: #A9ADB2;
		color: black;
	`
	activeButtonStyle = `
		border-color: #A9ADB2;
		box-shadow: 0 2px 6px rgba(31, 31, 38, 0);
		transform: translateY(1px);
	`
	focusActiveButtonStyle = `
		outline: 0;
	`
)

func oAuthButton(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("a", map[string]string{}, children).Style(baseButtonStyle+`
			display: flex;
			margin-left: 0;
			margin-right: 0;
			width: 100%;

			background: #323940;
			background-image: linear-gradient(-25deg, black, rgba(0,0,0,0));
			background-size: 300%;
			color: white;
			border-color: black;
		`).SelectorStyle(":focus", focusHoverButtonStyle+focusActiveButtonStyle).
			SelectorStyle(":hover", focusHoverButtonStyle).
			SelectorStyle(":active", activeButtonStyle+focusActiveButtonStyle).
			SelectorStyle(":hover", `
				background-color: #A9ADB2;
				color: white;
			`),
	}
}

func button(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("button", map[string]string{}, children).Style(baseButtonStyle+`
			display: flex;
			margin-left: 0;
			margin-right: 0;
			width: 100%;
		`).SelectorStyle(":focus", focusHoverButtonStyle+focusActiveButtonStyle).
			SelectorStyle(":hover", focusHoverButtonStyle).
			SelectorStyle(":active", activeButtonStyle+focusActiveButtonStyle),
	}
}

func inputLabel(attrs map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("label", attrs, children).Style(`
			display: block;
			padding-bottom: 2px;
			font-size: 14px;
			cursor: pointer;
		`).SelectorStyle(" > span", `
			color: #575B5F;
		`),
	}
}

func input(attrs map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("input", attrs, children).Style(`
			appearance: none;
			background: white;
			border-radius: 4px;
			border: 1px solid #C4C7CA;
			font-size: 16px;
			font-variant-ligatures: none;
			height: 2.5rem;
			padding: 0.5rem 0.75rem;
			width: 100%;
		`).SelectorStyle(":hover", `
			border-color: #81878e;
		`).SelectorStyle(":focus", `
			border-color: #81878e;
			box-shadow: 0 0 0 2px rgba(87,100,198,0.15);
			outline: 0;
		`).SelectorStyle("[type='range']", `
			padding: 0;
			margin: 1rem;
			cursor: ew-resize;
		`).SelectorStyle("[disabled]", `
			background: #ECECED;
			border-color: #ECECED;
		`),
	}
}

func text(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("p", map[string]string{}, children).Style(`
			line-height: 1.3;
			font-size: 14px;
			color: #575B5F;
			text-align: center
		`).SelectorStyle(" strong", `
			font-weight: 500;
		`).SelectorStyle(" a", `
			color: black;
			font-weight: 500;
		`).SelectorStyle(" a:hover", `
			color: #1E50DA;
		`).SelectorStyle(" a:active", `
			display: inline-block;
			text-decoration: none;
			transform: translateY(1px);
		`),
	}
}
