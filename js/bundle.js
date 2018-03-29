(function () {
    const loc = window.location
    const isSecure = loc.protocol === 'https:'

    const reload = document.createElement('div')
    reload.setAttribute('class', 'reloader')
    const link = document.createElement('a')
    link.setAttribute('href', loc.href)
    link.appendChild(document.createTextNode("New message on this thread, click to reload"))
    reload.appendChild(link)

    function checkUpdate(data) {
        const rid = data.record.toString()
        const recs = document.querySelectorAll("[data-record]")
        for (let i = 0; i < recs.length; i++) {
            const e = recs[i];
            const id = e.getAttribute('data-record')
            if (rid === id) {
                // const h = document.querySelector('.header')
                document.body.appendChild(reload)
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
})()