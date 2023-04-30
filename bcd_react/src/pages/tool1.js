import {Alert, FormControl, Grid, InputLabel, Paper, Select} from "@mui/material";
import React, {useState} from "react";
import axios from "axios";

const support_file_type = ['image/png', 'image/jpg', 'image/jpeg']

const select_language_baidu = [
    {value: 'CHN_ENG', viewValue: '中英文'},
    {value: 'ENG', viewValue: '英语'},
    {value: 'RUS', viewValue: '俄语'},
    {value: 'JAP', viewValue: '日语'},
    {value: 'KOR', viewValue: '韩语'},
    {value: 'FRE', viewValue: '法语'},
];

const language_to_baidu = {
    'CHN_ENG': 'zh',
    'ENG': 'en',
    'RUS': 'ru',
    'JAP': 'jp',
    'KOR': 'kor',
    'FRE': 'fra',
}

function Tool1() {

    const [selectVal, setSelectVal] = useState('CHN_ENG')
    const [imgSrc, setImgSrc] = useState('')
    const [tip1, setTip1] = useState('')
    const [tip2, setTip2] = useState('')
    const [res1, setRes1] = useState('')
    const [res2, setRes2] = useState('')

    function onDivPaste(event) {
        //首先显示图片
        let clipboardData = event.clipboardData;
        if (clipboardData != null && clipboardData.types[0] === 'Files') {
            let file = clipboardData.files.item(0);
            let fileType = '';
            if (file != null) {
                fileType = file.type
                console.log(fileType)
                if (support_file_type.indexOf(fileType.toLowerCase()) !== -1) {
                    triggerOcrText_baidu(file)
                } else {
                    console.log('file type[' + fileType + ']not support')
                }
            }

        }
    }

    function triggerOcrText_baidu(file) {
        let fr = new FileReader();//实例FileReader对象
        fr.readAsDataURL(file);//把上传的文件对象转换成url
        fr.onload = function (e) {
            console.log(e);
            let url = this.result.toString();//上传文件的URL
            // console.log(url)
            let base64 = url.split('base64,')[1]
            // console.log(base64)
            setImgSrc(url)
            setTip1('正在识别中')

            axios.post('/api/ocr/image/baiduAI', selectVal + ',' + base64, {
                headers: {"content-type": "text/html;charset=utf-8"},
            }).then(resp => {
                console.log(resp)
                let d = resp.data;
                if (d.code===0){
                    setTip1('有结果了')
                    setRes1(d.data)
                }else{
                    setTip1(d.message)
                }
            }, err => {
                console.log(err)
                setTip1(err)
            })

            setTip2('正在识别中')
            axios.post('/api/ocr/image/baiduFanyi', language_to_baidu[selectVal] + ',' + base64, {
                headers: {"content-type": "text/html;charset=utf-8"},
            }).then(resp => {
                console.log(resp)
                let d = resp.data;
                if (d.code===0){
                    setTip2('有结果了')
                    setRes2(d.data)
                }else{
                    setTip2(d.message)
                }
            }, err => {
                console.log(err)
                setTip2(err)
            })

        }
    }

    function onSelectChange(e) {
        setSelectVal(e.target.value)
    }

    return (
        <div onPaste={onDivPaste}>
            <div >
                <FormControl sx={{m: 1, minWidth: 120}}>
                    <InputLabel htmlFor="s">语言选择</InputLabel>
                    <Select native id="s" onChange={onSelectChange} label="语言选择" value={selectVal}>
                        {
                            select_language_baidu.map(item => {
                                return (
                                    <option key={`language_select_${item.value}`}
                                            value={item.value}>{item.viewValue}</option>
                                )
                            })
                        }
                    </Select>
                </FormControl>
            </div>

            <Grid container spacing={1}>
                <Grid item xs={6} style={{height: 600}}>
                    <Paper style={{width: "100%", height: "100%",backgroundColor:"lightskyblue"}}>
                        <img style={{marginTop:10,marginLeft:10}} alt='' width="95%" src={imgSrc}/>
                    </Paper>
                </Grid>
                <Grid item container xs={6} direction='column' spacing={1} style={{height: 600}}>
                    <Grid xs={6}  item>
                        <Paper style={{width: "100%", height: "100%"}}>
                            <Alert severity="info">{tip1}</Alert>
                            {res1}
                        </Paper>
                    </Grid>
                    <Grid xs={6} item>
                        <Paper style={{width: "100%", height: "100%"}}>
                            <Alert severity="info">{tip2}</Alert>
                            {res2}
                        </Paper>
                    </Grid>
                </Grid>
            </Grid>
        </div>
    );
}


export default Tool1