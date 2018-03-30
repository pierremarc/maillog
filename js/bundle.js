/*
 *  Copyright (C) 2018 Pierre Marchand <pierre.m@atelier-cartographique.be>
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as published by
 *  the Free Software Foundation, version 3 of the License.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

function start() {
    const loc = window.location
    const isSecure = loc.protocol === 'https:'

    const page = document.body.getAttribute("data-page")


    function reloader(msg) {
        const reload = document.createElement('div')
        reload.setAttribute('class', 'reloader')
        const reloadInner = document.createElement('div')
        reloadInner.setAttribute('class', 'reloader-inner')
        const link = document.createElement('a')
        link.setAttribute('href', loc.href)
        link.appendChild(document.createTextNode(msg))
        reloadInner.appendChild(link)
        reload.appendChild(reloadInner)
        document.body.appendChild(reload)
    }

    function checkUpdate(data) {
        if ('thread' === page) {
            const recs = document.querySelectorAll("[data-topic]")
            for (let i = 0; i < recs.length; i++) {
                const e = recs[i];
                const topic = e.getAttribute('data-topic')
                if (topic === data.topic) {
                    reloader(`Ther's a new message in ${data.topic}, click to reload`)
                }
            }
        }
        else if ('message' === page) {
            const rid = data.parent.toString()
            const recs = document.querySelectorAll("[data-record]")
            for (let i = 0; i < recs.length; i++) {
                const e = recs[i];
                const id = e.getAttribute('data-record')
                if (rid === id) {
                    reloader(`There's a new message in this thread, click to reload`)
                    const previousClass = e.getAttribute('class')
                    e.setAttribute('class', previousClass + ' new-reply')
                }
            }
        }
    }


    function connect() {
        const host = loc.host
        const socket = new WebSocket(`${isSecure ? 'wss' : 'ws'}://${host}/.notifications`)

        socket.addEventListener('open', function (event) {
            socket.send('Hello Server!');
        });

        socket.addEventListener('message', function (event) {
            console.log('Message from server ', event.data);
            try {
                const data = JSON.parse(event.data)
                checkUpdate(data)
            }
            catch (err) {
                console.error('Error processing data ', event.data);
            }
        });
    }


    connect()
}


document.onreadystatechange = function () {
    if ('interactive' === document.readyState) {
        start();
    }
};