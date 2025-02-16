import http from 'k6/http';
import { check } from 'k6';

const url = 'http://localhost:8080/api/auth';


export const options = {
    stages: [
        { duration: '5s', target: 700 }
    ],
};

export default function () {
    let data = { username: 'Bert', password:'123'};


    let res = http.post(url, JSON.stringify(data), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(res, { 'status was 200': (r) => r.status == 200 });

}
