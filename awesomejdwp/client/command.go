package client

import "fmt"

// cmdSet is the namespace for a command identifier.
type cmdSet uint8

// cmdId is a command in a command Set.
type cmdId uint8

type Command struct {
	Set cmdSet
	Id  cmdId
}

func (c Command) String() string {
	return fmt.Sprintf("%v.%v", c.Set, cmdNames[c])
}

const (
	cmdSetVirtualMachine       = cmdSet(1)
	cmdSetReferenceType        = cmdSet(2)
	cmdSetClassType            = cmdSet(3)
	cmdSetArrayType            = cmdSet(4)
	cmdSetInterfaceType        = cmdSet(5)
	cmdSetMethod               = cmdSet(6)
	cmdSetField                = cmdSet(8)
	cmdSetObjectReference      = cmdSet(9)
	cmdSetStringReference      = cmdSet(10)
	cmdSetThreadReference      = cmdSet(11)
	cmdSetThreadGroupReference = cmdSet(12)
	cmdSetArrayReference       = cmdSet(13)
	cmdSetClassLoaderReference = cmdSet(14)
	cmdSetEventRequest         = cmdSet(15)
	cmdSetStackFrame           = cmdSet(16)
	cmdSetClassObjectReference = cmdSet(17)
	cmdSetEvent                = cmdSet(64)
)

func (c cmdSet) String() string {
	switch c {
	case cmdSetVirtualMachine:
		return "VirtualMachine"
	case cmdSetReferenceType:
		return "ReferenceType"
	case cmdSetClassType:
		return "ClassType"
	case cmdSetArrayType:
		return "ArrayType"
	case cmdSetInterfaceType:
		return "InterfaceType"
	case cmdSetMethod:
		return "Method"
	case cmdSetField:
		return "Field"
	case cmdSetObjectReference:
		return "ObjectReference"
	case cmdSetStringReference:
		return "StringReference"
	case cmdSetThreadReference:
		return "ThreadReference"
	case cmdSetThreadGroupReference:
		return "ThreadGroupReference"
	case cmdSetArrayReference:
		return "ArrayReference"
	case cmdSetClassLoaderReference:
		return "ClassLoaderReference"
	case cmdSetEventRequest:
		return "EventRequest"
	case cmdSetStackFrame:
		return "StackFrame"
	case cmdSetClassObjectReference:
		return "ClassObjectReference"
	case cmdSetEvent:
		return "Event"
	}
	return fmt.Sprint(int(c))
}

var (
	cmdVirtualMachineVersion               = Command{cmdSetVirtualMachine, 1}
	cmdVirtualMachineClassesBySignature    = Command{cmdSetVirtualMachine, 2}
	cmdVirtualMachineAllClasses            = Command{cmdSetVirtualMachine, 3}
	cmdVirtualMachineAllThreads            = Command{cmdSetVirtualMachine, 4}
	cmdVirtualMachineTopLevelThreadGroups  = Command{cmdSetVirtualMachine, 5}
	cmdVirtualMachineDispose               = Command{cmdSetVirtualMachine, 6}
	cmdVirtualMachineIDSizes               = Command{cmdSetVirtualMachine, 7}
	cmdVirtualMachineSuspend               = Command{cmdSetVirtualMachine, 8}
	cmdVirtualMachineResume                = Command{cmdSetVirtualMachine, 9}
	cmdVirtualMachineExit                  = Command{cmdSetVirtualMachine, 10}
	cmdVirtualMachineCreateString          = Command{cmdSetVirtualMachine, 11}
	cmdVirtualMachineCapabilities          = Command{cmdSetVirtualMachine, 12}
	cmdVirtualMachineClassPaths            = Command{cmdSetVirtualMachine, 13}
	cmdVirtualMachineDisposeObjects        = Command{cmdSetVirtualMachine, 14}
	cmdVirtualMachineHoldEvents            = Command{cmdSetVirtualMachine, 15}
	cmdVirtualMachineReleaseEvents         = Command{cmdSetVirtualMachine, 16}
	cmdVirtualMachineCapabilitiesNew       = Command{cmdSetVirtualMachine, 17}
	cmdVirtualMachineRedefineClasses       = Command{cmdSetVirtualMachine, 18}
	cmdVirtualMachineSetDefaultStratum     = Command{cmdSetVirtualMachine, 19}
	cmdVirtualMachineAllClassesWithGeneric = Command{cmdSetVirtualMachine, 20}

	cmdReferenceTypeSignature            = Command{cmdSetReferenceType, 1}
	cmdReferenceTypeClassLoader          = Command{cmdSetReferenceType, 2}
	cmdReferenceTypeModifiers            = Command{cmdSetReferenceType, 3}
	cmdReferenceTypeFields               = Command{cmdSetReferenceType, 4}
	cmdReferenceTypeMethods              = Command{cmdSetReferenceType, 5}
	cmdReferenceTypeGetValues            = Command{cmdSetReferenceType, 6}
	cmdReferenceTypeSourceFile           = Command{cmdSetReferenceType, 7}
	cmdReferenceTypeNestedTypes          = Command{cmdSetReferenceType, 8}
	cmdReferenceTypeStatus               = Command{cmdSetReferenceType, 9}
	cmdReferenceTypeInterfaces           = Command{cmdSetReferenceType, 10}
	cmdReferenceTypeClassObject          = Command{cmdSetReferenceType, 11}
	cmdReferenceTypeSourceDebugExtension = Command{cmdSetReferenceType, 12}
	cmdReferenceTypeSignatureWithGeneric = Command{cmdSetReferenceType, 13}
	cmdReferenceTypeFieldsWithGeneric    = Command{cmdSetReferenceType, 14}
	cmdReferenceTypeMethodsWithGeneric   = Command{cmdSetReferenceType, 15}

	cmdClassTypeSuperclass   = Command{cmdSetClassType, 1}
	cmdClassTypeSetValues    = Command{cmdSetClassType, 2}
	cmdClassTypeInvokeMethod = Command{cmdSetClassType, 3}
	cmdClassTypeNewInstance  = Command{cmdSetClassType, 4}

	cmdArrayTypeNewInstance = Command{cmdSetArrayType, 1}

	cmdMethodTypeLineTable                = Command{cmdSetMethod, 1}
	cmdMethodTypeVariableTable            = Command{cmdSetMethod, 2}
	cmdMethodTypeBytecodes                = Command{cmdSetMethod, 3}
	cmdMethodTypeIsObsolete               = Command{cmdSetMethod, 4}
	cmdMethodTypeVariableTableWithGeneric = Command{cmdSetMethod, 5}

	cmdObjectReferenceReferenceType     = Command{cmdSetObjectReference, 1}
	cmdObjectReferenceGetValues         = Command{cmdSetObjectReference, 2}
	cmdObjectReferenceSetValues         = Command{cmdSetObjectReference, 3}
	cmdObjectReferenceMonitorInfo       = Command{cmdSetObjectReference, 5}
	cmdObjectReferenceInvokeMethod      = Command{cmdSetObjectReference, 6}
	cmdObjectReferenceDisableCollection = Command{cmdSetObjectReference, 7}
	cmdObjectReferenceEnableCollection  = Command{cmdSetObjectReference, 8}
	cmdObjectReferenceIsCollected       = Command{cmdSetObjectReference, 9}

	cmdStringReferenceValue = Command{cmdSetStringReference, 1}

	cmdThreadReferenceName                    = Command{cmdSetThreadReference, 1}
	cmdThreadReferenceSuspend                 = Command{cmdSetThreadReference, 2}
	cmdThreadReferenceResume                  = Command{cmdSetThreadReference, 3}
	cmdThreadReferenceStatus                  = Command{cmdSetThreadReference, 4}
	cmdThreadReferenceThreadGroup             = Command{cmdSetThreadReference, 5}
	cmdThreadReferenceFrames                  = Command{cmdSetThreadReference, 6}
	cmdThreadReferenceFrameCount              = Command{cmdSetThreadReference, 7}
	cmdThreadReferenceOwnedMonitors           = Command{cmdSetThreadReference, 8}
	cmdThreadReferenceCurrentContendedMonitor = Command{cmdSetThreadReference, 9}
	cmdThreadReferenceStop                    = Command{cmdSetThreadReference, 10}
	cmdThreadReferenceInterrupt               = Command{cmdSetThreadReference, 11}
	cmdThreadReferenceSuspendCount            = Command{cmdSetThreadReference, 12}

	cmdThreadGroupReferenceName     = Command{cmdSetThreadGroupReference, 1}
	cmdThreadGroupReferenceParent   = Command{cmdSetThreadGroupReference, 2}
	cmdThreadGroupReferenceChildren = Command{cmdSetThreadGroupReference, 3}

	cmdArrayReferenceLength    = Command{cmdSetArrayReference, 1}
	cmdArrayReferenceGetValues = Command{cmdSetArrayReference, 2}
	cmdArrayReferenceSetValues = Command{cmdSetArrayReference, 3}

	cmdClassLoaderReferenceVisibleClasses = Command{cmdSetClassLoaderReference, 1}

	cmdEventRequestSet                 = Command{cmdSetEventRequest, 1}
	cmdEventRequestClear               = Command{cmdSetEventRequest, 2}
	cmdEventRequestClearAllBreakpoints = Command{cmdSetEventRequest, 3}

	cmdStackFrameGetValues  = Command{cmdSetStackFrame, 1}
	cmdStackFrameSetValues  = Command{cmdSetStackFrame, 2}
	cmdStackFrameThisObject = Command{cmdSetStackFrame, 3}
	cmdStackFramePopFrames  = Command{cmdSetStackFrame, 4}

	cmdClassObjectReferenceReflectedType = Command{cmdSetClassObjectReference, 1}

	cmdEventComposite = Command{cmdSetEvent, 1}
)

var cmdNames = map[Command]string{}

func init() {
	register := func(c Command, n string) {
		if _, e := cmdNames[c]; e {
			panic("command already registered")
		}
		cmdNames[c] = n
	}
	register(cmdVirtualMachineVersion, "Version")
	register(cmdVirtualMachineClassesBySignature, "ClassesBySignature")
	register(cmdVirtualMachineAllClasses, "AllClasses")
	register(cmdVirtualMachineAllThreads, "AllThreads")
	register(cmdVirtualMachineTopLevelThreadGroups, "TopLevelThreadGroups")
	register(cmdVirtualMachineDispose, "Dispose")
	register(cmdVirtualMachineIDSizes, "IDSizes")
	register(cmdVirtualMachineSuspend, "Suspend")
	register(cmdVirtualMachineResume, "Resume")
	register(cmdVirtualMachineExit, "Exit")
	register(cmdVirtualMachineCreateString, "CreateString")
	register(cmdVirtualMachineCapabilities, "Capabilities")
	register(cmdVirtualMachineClassPaths, "ClassPaths")
	register(cmdVirtualMachineDisposeObjects, "DisposeObjects")
	register(cmdVirtualMachineHoldEvents, "HoldEvents")
	register(cmdVirtualMachineReleaseEvents, "ReleaseEvents")
	register(cmdVirtualMachineCapabilitiesNew, "CapabilitiesNew")
	register(cmdVirtualMachineRedefineClasses, "RedefineClasses")
	register(cmdVirtualMachineSetDefaultStratum, "SetDefaultStratum")
	register(cmdVirtualMachineAllClassesWithGeneric, "AllClassesWithGeneric")

	register(cmdReferenceTypeSignature, "Signature")
	register(cmdReferenceTypeClassLoader, "ClassLoader")
	register(cmdReferenceTypeModifiers, "Modifiers")
	register(cmdReferenceTypeFields, "Fields")
	register(cmdReferenceTypeMethods, "Methods")
	register(cmdReferenceTypeGetValues, "GetValues")
	register(cmdReferenceTypeSourceFile, "SourceFile")
	register(cmdReferenceTypeNestedTypes, "NestedTypes")
	register(cmdReferenceTypeStatus, "Status")
	register(cmdReferenceTypeInterfaces, "Interfaces")
	register(cmdReferenceTypeClassObject, "ClassObject")
	register(cmdReferenceTypeSourceDebugExtension, "SourceDebugExtension")
	register(cmdReferenceTypeSignatureWithGeneric, "SignatureWithGeneric")
	register(cmdReferenceTypeFieldsWithGeneric, "FieldsWithGeneric")
	register(cmdReferenceTypeMethodsWithGeneric, "MethodsWithGeneric")

	register(cmdClassTypeSuperclass, "Superclass")
	register(cmdClassTypeSetValues, "SetValues")
	register(cmdClassTypeInvokeMethod, "InvokeMethod")
	register(cmdClassTypeNewInstance, "NewInstance")

	register(cmdArrayTypeNewInstance, "NewInstance")

	register(cmdMethodTypeLineTable, "LineTable")
	register(cmdMethodTypeVariableTable, "VariableTable")
	register(cmdMethodTypeBytecodes, "Bytecodes")
	register(cmdMethodTypeIsObsolete, "IsObsolete")
	register(cmdMethodTypeVariableTableWithGeneric, "VariableTableWithGeneric")

	register(cmdObjectReferenceReferenceType, "ReferenceType")
	register(cmdObjectReferenceGetValues, "GetValues")
	register(cmdObjectReferenceSetValues, "SetValues")
	register(cmdObjectReferenceMonitorInfo, "MonitorInfo")
	register(cmdObjectReferenceInvokeMethod, "InvokeMethod")
	register(cmdObjectReferenceDisableCollection, "DisableCollection")
	register(cmdObjectReferenceEnableCollection, "EnableCollection")
	register(cmdObjectReferenceIsCollected, "IsCollected")

	register(cmdStringReferenceValue, "Value")

	register(cmdThreadReferenceName, "Name")
	register(cmdThreadReferenceSuspend, "Suspend")
	register(cmdThreadReferenceResume, "Resume")
	register(cmdThreadReferenceStatus, "Status")
	register(cmdThreadReferenceThreadGroup, "ThreadGroup")
	register(cmdThreadReferenceFrames, "Frames")
	register(cmdThreadReferenceFrameCount, "FrameCount")
	register(cmdThreadReferenceOwnedMonitors, "OwnedMonitors")
	register(cmdThreadReferenceCurrentContendedMonitor, "CurrentContendedMonitor")
	register(cmdThreadReferenceStop, "Stop")
	register(cmdThreadReferenceInterrupt, "Interrupt")
	register(cmdThreadReferenceSuspendCount, "SuspendCount")

	register(cmdThreadGroupReferenceName, "Name")
	register(cmdThreadGroupReferenceParent, "Parent")
	register(cmdThreadGroupReferenceChildren, "Children")

	register(cmdArrayReferenceLength, "Length")
	register(cmdArrayReferenceGetValues, "GetValues")
	register(cmdArrayReferenceSetValues, "SetValues")

	register(cmdClassLoaderReferenceVisibleClasses, "VisibleClasses")

	register(cmdEventRequestSet, "Set")
	register(cmdEventRequestClear, "Clear")
	register(cmdEventRequestClearAllBreakpoints, "ClearAllBreakpoints")

	register(cmdStackFrameGetValues, "GetValues")
	register(cmdStackFrameSetValues, "SetValues")
	register(cmdStackFrameThisObject, "ThisObject")
	register(cmdStackFramePopFrames, "PopFrames")

	register(cmdClassObjectReferenceReflectedType, "ReflectedType")

	register(cmdEventComposite, "Composite")
}
