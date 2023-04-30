import React, {useState} from 'react';
import {Alert, Button, Grid, Paper, TextField} from "@mui/material";
import {useNavigate} from "react-router-dom";
import axios from "axios";

function Login(props) {
    const [username_val, set_username_val] = useState('')
    const [password_val, set_password_val] = useState('')
    const [alert_data, set_alert_data] = useState({
        info: '',
        severity: 'info'
    })

    let navigate = useNavigate();

    function changeAlert(info, severity) {
        set_alert_data({
            info: info,
            severity: severity
        })
    }

    function validate() {
        let msg = ''
        if (username_val.trim() === '') {
            msg += '用户名不能为空'
        }
        if (password_val.trim() === '') {
            if (msg !== '') {
                msg += ','
            }
            msg += '密码不能为空'
        }
        if (msg !== '') {
            changeAlert(msg, 'error')
            return false
        } else {
            return true
        }
    }

    function login() {
        changeAlert('验证中', 'info')
        if (validate()) {
            changeAlert('验证通过', 'info')
            axios.post("/api/user/login", {
                username: username_val,
                password: password_val,
            }, {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }).then(
                res => {
                    console.log(res)
                    let d = res.data;
                    if (d['code']===0) {
                        changeAlert('登陆成功', 'success')
                        navigate('/main/photo')
                    } else {
                        changeAlert(`登陆失败[${d['message']}]`, 'error')
                    }
                }
                , err => {
                    console.log(err)
                    changeAlert(`请求失败[${err}]`, 'error')
                }
            );
        }
    }

    return (
        <Paper sx={{
            p: 2,
            margin: 'auto',
            maxWidth: 500,
            flexGrow: 1,
            marginTop: '200px',
            paddingTop: '50px',
            backgroundColor: '#b4e7b5'
        }}>
            <Grid container
                  direction="column"
                  spacing={2}>
                <Grid item textAlign="center">
                    <TextField
                        label="用户名"
                        value={username_val}
                        onChange={e => set_username_val(e.target.value)}
                        variant="outlined"/>
                </Grid>

                <Grid item textAlign="center">
                    <TextField
                        type="password"
                        value={password_val}
                        onChange={e => set_password_val(e.target.value)}
                        label="密码" variant="outlined"/>
                </Grid>

                <Grid item textAlign="center">
                    <Button onClick={login} variant="outlined">登陆</Button>
                    <Button style={{marginLeft: 20}} onClick={login} variant="outlined">重置</Button>
                </Grid>

                <Grid item textAlign="center">
                    <Alert sx={{display: alert_data.info === '' ? 'none' : ''}} variant="outlined"
                           severity={alert_data.severity}>
                        {alert_data.info}
                    </Alert>
                </Grid>
            </Grid>
        </Paper>
    );
}

export default Login;
