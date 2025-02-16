import http from 'k6/http';
import { check } from 'k6';

const url = 'http://localhost:8080/api/auth';


export const options = {
    stages: [
        { duration: '5s', target: 10 }
    ],
};

export default function () {
    let data = { username: 'Bert', password:'123'};

    // Using a JSON string as body
    let res = http.post(url, JSON.stringify(data), {
        headers: { 'Content-Type': 'application/json' },
    });
   // console.log(res.json()); // Bert
    check(res, { 'status was 200': (r) => r.status == 200 });
    // Using an object as body, the headers will automatically include
    // 'Content-Type: application/x-www-form-urlencoded'.
    //res = http.post(url, data);
    //console.log(res.json().form.name); // Bert

    // Using a binary array as body. Make sure to open() the file as binary
    // (with the 'b' argument).
   // http.post(url, logoBin, { headers: { 'Content-Type': 'image/png' } });

    // Using an ArrayBuffer as body. Make sure to pass the underlying ArrayBuffer
    // instance to http.post(), and not the TypedArray view.
   // data = new Uint8Array([104, 101, 108, 108, 111]);
   // http.post(url, data.buffer, { headers: { 'Content-Type': 'image/png' } });
}
