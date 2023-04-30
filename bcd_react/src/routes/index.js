import {Navigate} from "react-router-dom";
import Login from "../pages/login";
import Main from "../pages/main";
import Photo from "../pages/photo";
import Video from "../pages/video";
import Tool from "../pages/tool";
import Tool1 from "../pages/tool1";
import Tool2 from "../pages/tool2";
import Tool3 from "../pages/tool3";


const routerTable = [
    {
        path: '/login',
        element: <Login/>
    },
    {
        path: '/main',
        element: <Main/>,
        children: [
            {
                path: 'photo',
                element: <Photo/>
            },
            {
                path: 'video',
                element: <Video/>
            },
            {
                path: 'tool',
                element: <Tool/>,
                children: [
                    {
                        path: 'tool1',
                        element: <Tool1/>,
                    },
                    {
                        path: 'tool2',
                        element: <Tool2/>,
                    },
                    {
                        path: 'tool3',
                        element: <Tool3/>,
                    }
                ]
            }
        ]
    },
    {
        path: '/',
        element: <Navigate to="/login"/>
    }
]

export default routerTable


