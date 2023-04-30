import React, {useEffect, useState} from 'react';
import {Backdrop, Container, ImageList, ImageListItem, Button, Alert, Grid, IconButton} from "@mui/material"
import {Delete} from "@mui/icons-material"
import axios from "axios";

function Photo(props) {
    const [open, setOpen] = React.useState(false);
    const [coverImage, setCoverImage] = React.useState('');
    const [photos, setPhotos] = React.useState([]);
    const [uploadAlert, setUploadAlert] = useState({
        info: '',
        severity: 'info'
    });
    const fileInput = React.useRef();

    const list=()=>{
        try {
            axios.get("/api/photo/list").then(res=>{
                const d = res.data;
                console.log(d)
                if (d['code'] === 0&&d['data']!=null) {
                    setPhotos(d['data'])
                }
            })
        } catch (err) {
            console.log(err)
        }
    };

    const handleClose = () => {
        setOpen(false);
    };
    const handleOpen = (img) => {
        setCoverImage(img)
        setOpen(!open);
    };



    const handleFileChange = (e) => {
        console.log(e.target.files)
        const formData = new FormData();
        formData.append("file", e.target.files[0])
        axios.post("/api/photo/upload", formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        }).then(
            res => {
                let d = res.data;
                if (d['code'] === 0) {
                    setUploadAlert({
                        info: '上传成功',
                        severity: 'success'
                    })
                    list()
                } else {
                    setUploadAlert({
                        info: `上传失败,[${d['message']}]`,
                        severity: 'error'
                    })
                }
                e.target.value=''
            }
            , err => {
                setUploadAlert({
                    info: `上传失败[${err}]`,
                    severity: 'error'
                })
                e.target.value=''
            }
        );
    };

    const handleUploadButtonClick = (e) => {
        fileInput.current.click()
    }

    const handleDelButtonClick = (e) => {
        const name = coverImage.substring(coverImage.lastIndexOf("/") + 1);
        axios.get("/api/photo/del", {
            params:{
                name: name
            }
        }).then(
            res => {
                let d = res.data;
                if (d['code'] === 0) {
                    setUploadAlert({
                        info: `删除图片[${name}]成功`,
                        severity: 'success'
                    })
                    list()
                } else {
                    setUploadAlert({
                        info: `删除图片[${name}]失败,[${d['message']}]`,
                        severity: 'error'
                    })
                }
            }
            , err => {
                setUploadAlert({
                    info: `删除图片[${name}]失败[${err}]`,
                    severity: 'error'
                })
            }
        );
    }

    //注册键盘监听
    useEffect(() => {
        const fn = e => {
            if (e.code === 'Space') {
                setOpen(old => {
                    if (old) {
                        return false
                    } else {
                        return old
                    }
                })
            }
            e.preventDefault()
        }
        window.addEventListener('keypress', fn)
        return () => {
            window.removeEventListener('keypress', fn)
        }
    }, [])

    //加载所有照片
    useEffect(() => {
        list()
    }, [])

    return (
        <Container>
            <Grid container>
                <Grid xs={1} item>
                    <Button onClick={handleUploadButtonClick} variant="contained">上传</Button>
                    <input style={{display: "none"}} ref={fileInput} type={"file"} onChange={handleFileChange}/>
                </Grid>
                <Grid xs={1} item>
                    <Button onClick={list} variant="contained">刷新</Button>
                </Grid>
                <Grid xs={10} item>
                    <Alert sx={{display: uploadAlert.info === '' ? 'none' : ''}} variant="contained"
                           severity={uploadAlert.severity}>
                        {uploadAlert.info}
                    </Alert>
                </Grid>
            </Grid>

            <ImageList sx={{maxHeight: 700}} cols={3}>
                {
                    photos.map(item => (
                            <ImageListItem sx={{":hover": {border: 4, borderColor: 'green'}}} key={item}>
                                <img
                                    src={`/api/photo/download/${item}`}
                                    alt=''
                                    loading="lazy"
                                    onClick={e => {
                                        handleOpen(`/api/photo/download/${item}`)
                                    }}
                                />
                            </ImageListItem>
                        )
                    )
                }
            </ImageList>
            <Backdrop
                sx={{color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1}}
                open={open}
                onClick={handleClose}
            >
                <Grid container>
                    <Grid xs={12} style={{textAlign: 'right'}} item>
                        <IconButton onClick={handleDelButtonClick} aria-label="delete" size="large" color="error">
                            <Delete/>
                        </IconButton>
                    </Grid>
                    <Grid xs={12} style={{textAlign: 'center'}} item>
                        <img height="700" src={coverImage} alt=''/>
                    </Grid>
                </Grid>
            </Backdrop>
        </Container>
    );
}

export default Photo;