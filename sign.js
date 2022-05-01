(async () => {
    const tools = {
        random_string: (len) => {
            let str = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
            let result = '';
            for (let i = 0; i < len; i++) {
                result += str[Math.floor(Math.random() * str.length)];
            }
            return result;
        },
    };
    const S = require('crypto');
    console.log(JSON.stringify(
        await ((
            ${SIGN}
        )(process.argv[2], JSON.parse(process.argv[3]), S))
    ));
})();