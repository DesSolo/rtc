import React from 'react';
import { Routes, Route } from 'react-router-dom';
import LayoutWithNav from "../../components/Layout/Layout.jsx";
import ProjectsList from "../../components/Projects/List.jsx";
import ConfigsList from "../../components/Configs/Detail.jsx";
import Releases from "../Releases/index.jsx";

const Index = () => {
    return (
        <Routes>
            {/* Маршруты с NavBar */}
            <Route path="/" element={<LayoutWithNav />}>
                <Route index element={<ProjectsList />} />
                <Route path="projects/:project/releases" element={<Releases />} />
                <Route path="projects/:project/envs/:environment/releases/:release/configs" element={<ConfigsList />} />
            </Route>

            {/* Маршруты без NavBar */}
            {/*<Route path="login" element={<Login />} />*/}
        </Routes>
    );
};

export default Index