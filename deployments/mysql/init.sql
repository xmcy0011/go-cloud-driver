-- 虚拟根节点，所有节点的祖先，闭包表无需插入，降低深度
insert into metadata(object_id, parent_id, name, object_type) values('00000000000000000000000000', '', 'url', '0');
-- 虚拟的根无需插入闭包表，增加额外的层级，且没多大用处
-- insert into metadata_closure(ancestor,descendant,depth) values('00000000000000000000000000','00000000000000000000000000',0);