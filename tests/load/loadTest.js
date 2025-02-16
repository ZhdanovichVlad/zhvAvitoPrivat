import http from 'k6/http';
import { check } from 'k6';


const authUrl = 'http://0.0.0.0:8080/api/auth';


const testUrl = 'http://0.0.0.0:8080/api/info';

export const options = {
    stages: [
        { duration: '5s', target: 850 },
        { duration: '5s', target: 1000 },
        { duration: '5s', target: 700 },
    ],
};

export default function () {

    let authData = { username: 'avito', password: '123' };
    let authRes = http.post(authUrl, JSON.stringify(authData), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(authRes, {
        'status was 200': (r) => r.status === 200,
        'received JWT token': (r) => r.json().token !== undefined,
    });


    let token = authRes.json().token;

    let headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`, // Добавляем токен в заголовок
    };


    let testRes = http.get(testUrl, { headers: headers });


    check(testRes, {
        'status was 200': (r) => r.status === 200,
    });
}