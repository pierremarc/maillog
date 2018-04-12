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
    /* width: 90%; */
    margin: 2em;
    position: relative;
}


body[data-page=root] > .header, body[data-page=thread] > .header {
    margin: 1em auto;
    padding-bottom: 1em;
    border-bottom: 0.5px solid black;
}

body[data-page=message] > .header {
    margin: auto auto 3em 12em;
    padding-bottom: 1em;
}

.header > h1.title {
    font-size: 180%;
    font-weight: bold;
}

.header > .bc-block {
margin-top: .3em;
}

.header > .bc-block > .bc {
    font-size: 80%;
    /* font-family: monospace;
    font-weight: bold; */
}

a.link {
    color: #c7002e;
    font-family: monospace;
    font-weight: bold;
}

.message-block, .answer-block {
    display: flex;
    position: relative;
    padding: 1em 1em 1em 0;
}

.message-header, .answer-header-block {
    flex: 0 0 12em;
    text-align: right;
    padding-right: 1em;
    line-height: 140%;
    font-size: 100%;
    font-family: monospace;
}

.section-link {
    text-decoration: none;
    font-weight: bold;
    font-size: 124%;
    background-color: black;
    color: black;
}

.section-link:hover {
    color:grey;
}

.message-body, .answer-body {
    padding: 0 0 0 1em;
    width: 42em;
}

.message-sender {
    /* font-style: italic; */
    font-size: 90%;
}

.message-date {
    font-size: 80%;
}

.answer-view {
    font-size: 80%;
}

 .message-header > .parent {
font-size: 80%;
}

/* .message-block > .message-header > a {
    font-size: 80%;
    position: absolute;
    right: 0;
    margin-top: 0.2em;
} */

 


p.paragraph {
    padding-bottom: 1em;
}

img {
    max-width:100%;
}

/* .answer {
    border-left: 0 solid #2196f3;
    background-color: white;
}

.answer.depth-1 { border-left-width: 1px; }
.answer.depth-2 { border-left-width: 2px; }
.answer.depth-3 { border-left-width: 4px; }
.answer.depth-4 { border-left-width: 6px;}
.answer.depth-5 { border-left-width: 8px;}
.answer.depth-6 {border-left-width: 10px;}
.answer.depth-7 { border-left-width: 12px;}
.answer.depth-8 { border-left-width: 14px;}
 */


.answer-link {
/* margin-top: .3em; */
}

.answer-link a.link {
    color:#2196f3;
    font-size: 80%;
}



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

.message-body,.answer-body p {
    line-height: 140%;
}

.reloader {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: auto;
    text-align: center;
}

.reloader-inner {
    background-color: rgb(243, 239, 239);
        padding: 12px;
}
.reloader a{
    font-size: 90%;
    font-weight: bold;
    color: #0d6903;
}


.topic-replyto-block {
    font-size: 80%;
    line-height: 120%;
    margin-bottom: 1em;
    border-bottom: 0.5px solid black;
    padding-bottom: 1em;
}

.first-subject {
    font-style: italic;
    color: black;
}
.first-subject:hover{
    color: grey;
}

.new-reply {
    background-color: #ffefbf;
}

@media screen and (orientation: portrait), screen and (max-width: 48em)  {
body {
    margin: 2em 1em;
}

body[data-page=message] > .header {
    margin: 1em auto;
    padding-bottom: 1em;
}

.message-block, .answer-block {
    display: initial;
    position: relative;
    padding: initial;
}

.message-header, .answer-header-block {
    flex: none;
    text-align: left;
    line-height: 140%;
    font-size: 100%;
    font-family: monospace;
}

.message-body, .answer-body {
    padding: 0.3em 0 .7em 0;
    width: 100%;
}

.message-sender {
    display: inline-block;
    font-size: 90%;
}

.message-date {
    display: inline-block;
    font-size: 80%;
}

.answer-view {
    display: inline-block;
    font-size: 80%;
}
.answer-link {
    display: inline-block;
}

 .message-header > .parent {
    display: inline-block;
font-size: 80%;
}

.answer {
    border:none;
    background-color: white;
}

}`

