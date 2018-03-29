package main


const CssReset = `/* http://meyerweb.com/eric/tools/css/reset/ 
   v2.0 | 20110126
   License: none (public domain)
*/

html, body, div, span, applet, object, iframe,
h1, h2, h3, h4, h5, h6, p, blockquote, pre,
a, abbr, acronym, address, big, cite, code,
del, dfn, em, img, ins, kbd, q, s, samp,
small, strike, strong, sub, sup, tt, var,
b, u, i, center,
dl, dt, dd, ol, ul, li,
fieldset, form, label, legend,
table, caption, tbody, tfoot, thead, tr, th, td,
article, aside, canvas, details, embed, 
figure, figcaption, footer, header, hgroup, 
menu, nav, output, ruby, section, summary,
time, mark, audio, video {
	margin: 0;
	padding: 0;
	border: 0;
	font-size: 100%;
	font: inherit;
	vertical-align: baseline;
}
/* HTML5 display-role reset for older browsers */
article, aside, details, figcaption, figure, 
footer, header, hgroup, menu, nav, section {
	display: block;
}
body {
	line-height: 1;
}
ol, ul {
	list-style: none;
}
blockquote, q {
	quotes: none;
}
blockquote:before, blockquote:after,
q:before, q:after {
	content: '';
	content: none;
}
table {
	border-collapse: collapse;
	border-spacing: 0;
}`

const CssStyle = `
body {
    font-family: sans-serif;
    width: 90%;
    margin: auto;
    position: relative;
}
@media screen and (max-width: 48em) {
    body {
        /* font-size:132%; */
}
}

body div.header {
    margin: 1em auto 1em auto;
    padding-bottom: 1em;
    border-bottom: 0.5px solid black;
}

body > div.header > h1.title {
        font-size: 140%;
    font-weight: bold;
}

body > div.header > div.bc-block {
        position: absolute;
    top: 0;
    right: 0;
}

body > div.header > div.bc-block > div.bc {
    font-size: 80%;
    /* font-family: monospace;
    font-weight: bold; */
}

a.link {
    color: #c7002e;
    font-family: monospace;
    font-weight: bold;
}

body > div.message-block {
    position: relative;
}

body > div.message-block > div.message-header {
    margin-bottom: 1em;
    padding-bottom: 1em;
    border-bottom: 0.5px solid black;
}

body > div.message-block > div.message-header > span.message-sender {
    font-style: italic;
    font-size: 90%;
}

body > div.message-block > div.message-header > a {
    font-size: 80%;
    position: absolute;
    right: 0;
    margin-top: 0.2em;
}

body > div.message-block > div.message-body {
    padding: 16px;
}

body > div.message-block > div.message-header > div.parent {
    position: absolute;
    top: -1em;
}
p.paragraph {
    padding-bottom: 1em;
}

img {
    max-width:100%;
}

body > div.answer {
border: 0.5px solid gray;
    padding: 1em;
}


div.answer-block {
    padding-bottom: 0.5em;
    margin-bottom: 0.5em;
    border-bottom: 0.5px solid gray;
}

div.answer-block > div.answer-header-block {
   position: relative;
    border-bottom: 0.5px solid gray;
    padding-bottom: 0.5em;
    margin-bottom: 0.5em;
}

div.answer-header-block > a.answer-link {
    position: absolute;
    top: 0;
    right: 0;
}

div.answer.depth-2 { margin-left: 1em; }
div.answer.depth-3 { margin-left: 2em; }
div.answer.depth-4 { margin-left: 3em; }
div.answer.depth-5 { margin-left: 4em; }
div.answer.depth-6 { margin-left: 5em; }
div.answer.depth-7 { margin-left: 6em; }
div.answer.depth-8 { margin-left: 7em; }


.topic {
    margin-bottom: 1em;
    font-size: 120%;
}
.topic-count{
    font-size: small;
}

span.topic-ts {
    font-size: small;
}
.message-item {
        margin-bottom: .5em;
}

.message-item-sender {
    font-size: 69%;
    font-style: italic;
}


span.message-item-ts {
    font-family: monospace;
    font-size: 64%;
    color: #9E9E9E;
}

p.body-par {
    line-height: 140%;
}`

