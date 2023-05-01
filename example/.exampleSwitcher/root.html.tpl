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
    <div style="display: flex">
        <nav>
            {{range .Examples}}
                <ul>
                    <li>
                        <a href="/{{ .Path }}">
                        {{ .Name }}
                        </a>
                    </li>
                </ul>
            {{end}}
        </nav>
    </div>
</body>

</html>