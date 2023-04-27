const axios = require('axios')
const fs = require('fs')

const nwords = ['zero', 'one', 'two', 'three', 'four', 'five', 'six', 'seven', 'eight', 'nine']
const ns = {'zero' : 0, 'one' : 1, 'two' : 2, 'three' : 3, 'four' : 4, 'five': 5, 'six': 6, 'seven': 7, 'eight': 8, 'nine': 9}
axios
    .get('https://quest.squadcast.tech/api/310620104093/worded_ip')
    .then(async function (res) {
        // console.log(res)
        let para = res.data.split('<div id="container">')[1].split('</div>')[0].trim();
        let nums = []
        buf = []
        for (word of para.split(" ")) {
            if (nwords.includes(word)) {
                buf.push(ns[word])
            } else if (word == 'point') {
                nums.push(buf.join(''))
                buf = []
            }
        }
        nums.push(buf.join('')) // final octet

        let ip = []
        let i = 0, l = 0;
        while (l < 4) {
            if (nums[i] > 255) {
                l = 0
                ip = []
            }
            else if (nums[i] <= 0 && l == 0) {
                l = 0
                ip = []
            }
            else if (nums[i] < 0) {
                l = 0
                ip = []
            }
            else 
                ip.push(nums[i])
            i++;
        }

        ans = ip.join('.')

        let config = {
            headers: {
                'Content-Type': 'text/plain'
            }
        }
        let code = fs.readFileSync('index.js', 'utf-8');
        await axios.post(`https://quest.squadcast.tech/api/310620104093/submit/worded_ip?ip=${ans}`, code, config)
    }       
)
