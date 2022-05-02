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
            async(t,e,S)=>{e=JSON.parse(JSON.stringify(e));e["private_key"]=t;let s=Object.keys(e).sort().map(t=>{return t+'='+e[t]}).join("&");let r=S.createHash("md5").update(s).digest().reduce((t,e)=>{return t+((e&255).toString(16).length===1?"0"+(e&255).toString(16):(e&255).toString(16))},"").toLowerCase();const a=t=>{let e=S.createHash("md5").update(t).digest("hex").toLowerCase();let s=e.substr(3,9);let r=e.substr(7,20);let a=e.substr(1,9);let o=t.split("").map(t=>{return t.charCodeAt(0).toString(16)}).join("");let l=tools.random_string(3);let n=[s,r,a,o].sort().join("")+l;let i=S.createHash("sha1").update(n).digest("hex").toLowerCase();let u=i+l;let d=o;let h=e.substr(3,26);let g=d+h;let p=g.split("").map(t=>{return t.charCodeAt(0).toString(16)}).join("").substr(1,17);let b=Buffer.from(u+p).toString("base64");let c=b+g.split("").map(t=>{return t.charCodeAt(0).toString(16)}).join("").substr(48,64)+"hexbase64";let m=S.createHash("md5").update(c).digest("hex").toLowerCase();return l+h+m};const o=t=>{let e=S.createHash("md5").update(t).digest("hex").toLowerCase();let s=e.substr(5,21);let r=e;let a=r.split("").map(t=>{return t.charCodeAt(0).toString(16)}).join("").substr(1,44);let o=S.createHash("md5").update(a).digest("hex").toLowerCase();let l=o.substr(5,16);let n=o.substr(0,16);let i=s+r+t;let u=S.createCipheriv("aes-128-cbc",l,n);let d=u.update(i,"utf8","base64");let h=u.final("base64");let g=d+h;let p=r.substr(3,9);let b=r.substr(7,20);let c=r.substr(8,9);let m=S.createHash("sha1").update([p,b,c,g].sort().join("")).digest("hex").toLowerCase()+tools.random_string(6);let C=tools.random_string(5);return C+m};return{sign:r,sesi:o(r),nars:a(r)}}
        )(process.argv[2], JSON.parse(process.argv[3]), S))
    ));
})();