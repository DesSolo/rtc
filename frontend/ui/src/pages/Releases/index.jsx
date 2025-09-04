import {useParams} from "react-router-dom";
import ReleasesList from "../../components/Releases/List.jsx";

const Releases = () => {
    const {project} = useParams();

    return (
        <ReleasesList project={project} />
    )
}

export default Releases;