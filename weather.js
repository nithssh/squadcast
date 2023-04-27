const axios = require('axios')
const fs = require('fs')

var promises = [];
promises.push(axios.get('https://quest.squadcast.tech/api/310620104093/weather/get?q=Hyderabad'))
promises.push(axios.get('https://quest.squadcast.tech/api/310620104093/weather/get?q=Bengaluru'))
Promise.all(promises).then(async function (ress) {
    let data = ress.map((x) => x.data)
    let code = fs.readFileSync('index.js', 'utf-8');
    
    // testcase
    let config = {
        headers: {
            'Content-Type': 'text/plain'
        }
    }
    try {
        if (data[0].wind.speed > data[1].wind.speed)
            await axios.post(`https://quest.squadcast.tech/api/310620104093/submit/weather?answer=${data[0].name}`, code, config)
        else
            await axios.post(`https://quest.squadcast.tech/api/310620104093/submit/weather?answer=${data[1].name}`, code, config)
    } catch (e) {
        console.error(e)
    }
})
