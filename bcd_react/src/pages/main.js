import React, {useEffect, useState} from 'react';
import {Outlet, useLocation, useNavigate} from "react-router-dom";
import {AppBar, Tab, Tabs} from "@mui/material";

function Main(props) {

    const [tabVal, setTabVal] = useState('photo')

    let navigate = useNavigate();
    let location = useLocation();

    useEffect(() => {
        let split = location.pathname.split("/");
        if (split.length > 2) {
            setTabVal(() => split[2])
        }else{
            navigate('photo')
        }
    }, [location.pathname, navigate])

    function handleChange(e, v) {
        navigate(v)
    }

    return (
        <div>
            <AppBar position="static">
                <Tabs
                    value={tabVal}
                    indicatorColor="secondary"
                    onChange={handleChange}
                    textColor="inherit"
                    variant="fullWidth"
                >
                    <Tab value='photo' label="照片"/>
                    <Tab value='video' label="视频"/>
                    <Tab value='tool' label="工具"/>
                </Tabs>
            </AppBar>
            <Outlet/>
        </div>
    );
}

export default Main;

