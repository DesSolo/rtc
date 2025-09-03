import { Layout, Menu } from "antd";

import {
  ProjectOutlined,
  AuditOutlined
} from '@ant-design/icons';

const { Sider } = Layout;

const Nav = () => {
    return (
        <Sider trigger={null} collapsible> // TODO: add collapse
        <div className="demo-logo-vertical" />
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

export default Nav;