import { Layout, Menu } from "antd";
import {
  ProjectOutlined,
  AuditOutlined
} from '@ant-design/icons';

const { Sider } = Layout;

const NavBar = ({collapsed}) => {
    return (
        <Sider trigger={null} collapsible collapsed={collapsed}>
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={['1']}
          items={[
            {
              key: '1',
              icon: <ProjectOutlined />,
              label: 'Projects',
            },
            {
              key: '2',
              icon: <AuditOutlined />,
              label: 'Audit logs',
            },
          ]}
        />
      </Sider>
    )
};

export default NavBar;