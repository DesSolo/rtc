import React, {useEffect} from 'react';
import { Routes, Route } from 'react-router-dom';
import LayoutWithNav from "../../components/Layout/Layout.jsx";
import ProjectsList from "../../components/Projects/List.jsx";
import ConfigsList from "../../components/Configs/Detail.jsx";
import Releases from "../Releases/index.jsx";
import Audit from "../Audit/index.jsx";
import Login from "../Login/index.jsx"
import Users from "../Users/index.jsx";


const Index = () => {
    useEffect(()=>{
        document.title = 'RTC'
    })
    return (
        <Routes>
            <Route path="/" element={<LayoutWithNav />}>
                <Route index element={<ProjectsList />} />
                <Route path="projects/:project/releases" element={<Releases />} />
                <Route path="projects/:project/envs/:environment/releases/:release/configs" element={<ConfigsList />} />
                <Route path="audit" element={<Audit />} />
                <Route path="/users" element={<Users />} />
            </Route>

            <Route path="login" element={<Login />} />
        </Routes>
    );
};

export default Index