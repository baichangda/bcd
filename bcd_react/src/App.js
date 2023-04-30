import './App.css';
import {useNavigate, useRoutes} from "react-router-dom";
import routeTable from './routes'
import axios from "axios";
import {useEffect} from "react";

function App() {
    const navigate = useNavigate();
    useEffect(() => {
        axios.interceptors.response.use(res => {
            if (res.status === 200 && res.headers['content-type'] === 'application/json' && res.data['code'] === 401) {
                navigate("/login", {
                    replace: true,
                    state: {
                        errorMsg:res.data['message']
                    }
                })
            }
            return res
        }, err => {
            return Promise.reject(err)
        })
    }, [navigate])

    let rt = useRoutes(routeTable);
    return (
        <>
            {rt}
        </>
    );
}

export default App;
