package classfile

import "fmt"

//  ClassFile {
//  	u4                 magic;
//  	u2                 minor_version;
//  	u2                 major_version;
//  	u2                 constant_pool_count;
//  	cp_info 		   constant_pool[constant_pool_count-1];
//  	u2                 access_flags;
//  	u2                 this_class;
//  	u2                 super_class;
//  	u2                 interfaces_count;
//  	u2                 interfaces[interfaces_count];
//  	u2                 fields_count;
//  	field_info         fields[fields_count];
//  	u2                 methods_count;
//  	method_info        methods[methods_count];
//  	u2                 attributes_count;
//  	attribute_info     attributes[attributes_count];
//  }

type ClassFile struct {
	//magic        uint32
	minorVersion uint16
	majorVersion uint16
	constantPool ConstantPool
	accessFlags  uint16
	thisClass    uint16
	superClass   uint16
	interfaces   []uint16
	fields       []*MemberInfo
	methods      []*MemberInfo
	attributes   []AttributeInfo
}

type MemberInfo struct {
	cp              ConstantPool
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributes      []AttributeInfo
}

type ConstantPool []ConstantInfo

type ConstantInfo interface {
	readInfo(reader *ClassReader)
}

type AttributeInfo interface {
	readInfo(reader *ClassReader)
}

// 函数把 []byte 解析成 ClassFile 结构体
func Parse(classData []byte) (classFile *ClassFile, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	classReader := &ClassReader{classData}
	classFile = &ClassFile{}
	classFile.read(classReader)
	return
}

// TODO
// ./jvm -Xjre "/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre" Test
// func (self *ClassFile) String() string {
// 	version := fmt.Sprintf("version: %v.%v\n", cf.MajorVersion(), cf.MinorVersion())
// 	constantsCount := fmt.Sprintf("constants count: %v\n", len(cf.ConstantPool()))
// 	accessFlags := fmt.Sprintf("access flags: 0x%x\n", cf.AccessFlags())
// 	thisClass := fmt.Sprintf("this class: %v\n", cf.ClassName())
// 	superClass := fmt.Sprintf("super class: %v\n", cf.SuperClassName())
// 	interfaces := fmt.Sprintf("interfaces: %v\n", cf.InterfaceNames())
// 	fieldsCount := fmt.Sprintf("fields count: %v\n", len(cf.Fields()))
// 	fields := ""
// 	for _, f := range cf.Fields() {
// 		fields += fmt.Sprintf("  %s\n", f.Name())
// 	}
// 	methodsCount := fmt.Sprintf("methods count: %v\n", len(cf.Methods()))
// 	methods := ""
// 	for _, m := range cf.Methods() {
// 		methods += fmt.Sprintf("  %s\n", m.Name())
// 	}
// 	return version + constantsCount + accessFlags + thisClass + superClass + interfaces + fieldsCount + fields + methodsCount + methods
// }
