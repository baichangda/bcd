import React, {useEffect, useState} from 'react';
import {Box, Tab, Tabs} from "@mui/material";
import {Outlet, useLocation, useNavigate} from "react-router-dom";

function Tool(props) {

    const [tabVal, setTabVal] = useState('tool1')

    let navigate = useNavigate();
    let location = useLocation();

    useEffect(() => {
        let split = location.pathname.split("/");
        console.log('------------', split)
        if (split.length > 3) {
            setTabVal(() => split[3])
        }else{
            navigate('tool1')
        }
    }, [location.pathname, navigate])

    function handleChange(e, v) {
        navigate(v)
    }

    return (
        <Box sx={{width: '100%'}}>
            <Tabs
                value={tabVal}
                onChange={handleChange}
                textColor="secondary"
                indicatorColor="secondary"
            >
                <Tab value='tool1' label="粘贴图片识别文字"/>
                <Tab value='tool2' label="粘贴图片识别表格(腾讯)"/>
                <Tab value='tool3' label="粘贴图片识别表格(百度)"/>
            </Tabs>
            <Outlet/>
        </Box>

    );
}

export default Tool;