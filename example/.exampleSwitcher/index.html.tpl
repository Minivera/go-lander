<!doctype html>
<!--
Copyright 2018 The Go Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
-->
<html lang="en">

<head>
	<meta charset="utf-8">
	<title>Go wasm</title>
<body>
	<!--
	Add the following polyfill for Microsoft Edge 17/18 support:
	<script src="https://cdn.jsdelivr.net/npm/text-encoding@0.7.0/lib/encoding.min.js"></script>
	(see https://caniuse.com/#feat=textencoder)
	-->
	<script src="index.js"></script>
	<script>
		if (!WebAssembly.instantiateStreaming) { // polyfill
			WebAssembly.instantiateStreaming = async (resp, importObject) => {
				const source = await (await resp).arrayBuffer();
				return await WebAssembly.instantiate(source, importObject);
			};
		}

		const go = new Go();
		let mod, inst;
		WebAssembly.instantiateStreaming(fetch("./main.wasm"), go.importObject).then((result) => {
			mod = result.module;
			inst = result.instance;

			run();
		}).catch((err) => {
			console.error(err);
		});

		async function run() {
			console.clear();
			await go.run(inst);
			inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
		}
	</script>

    <div style="display: flex">
        <nav>
            {{range .Examples}}
                <ul>
                    <li>
                        <a href="/{{ .Path }}">
                        {{ if .Active }}<b>{{end}}{{ .Name }}{{ if .Active }}</b>{{end}}
                        </a>
                    </li>
                </ul>
            {{end}}
        </nav>
        <div style="margin-left: 3rem; display: flex; flex-direction: column;">
	        <div id="app"></div>
	        <div style="margin-top: 1rem;">
	            Find this example's code at <a target="_blank" href="{{ .Url }}">{{ .Url }}</a>.
	        </div>
        </div>
    </div>
</body>

</html>