package client

import "fmt"

// CmdSet is the namespace for a command identifier.
type CmdSet uint8

// CmdId is a command in a command Set.
type CmdId uint8

type Command struct {
	Set CmdSet
	Id  CmdId
}

func (c Command) String() string {
	return fmt.Sprintf("%v.%v", c.Set, CmdNames[c])
}

const (
	CmdSetVirtualMachine       = CmdSet(1)
	CmdSetReferenceType        = CmdSet(2)
	CmdSetClassType            = CmdSet(3)
	CmdSetArrayType            = CmdSet(4)
	CmdSetInterfaceType        = CmdSet(5)
	CmdSetMethod               = CmdSet(6)
	CmdSetField                = CmdSet(8)
	CmdSetObjectReference      = CmdSet(9)
	CmdSetStringReference      = CmdSet(10)
	CmdSetThreadReference      = CmdSet(11)
	CmdSetThreadGroupReference = CmdSet(12)
	CmdSetArrayReference       = CmdSet(13)
	CmdSetClassLoaderReference = CmdSet(14)
	CmdSetEventRequest         = CmdSet(15)
	CmdSetStackFrame           = CmdSet(16)
	CmdSetClassObjectReference = CmdSet(17)
	CmdSetEvent                = CmdSet(64)
)

func (c CmdSet) String() string {
	switch c {
	case CmdSetVirtualMachine:
		return "VirtualMachine"
	case CmdSetReferenceType:
		return "ReferenceType"
	case CmdSetClassType:
		return "ClassType"
	case CmdSetArrayType:
		return "ArrayType"
	case CmdSetInterfaceType:
		return "InterfaceType"
	case CmdSetMethod:
		return "Method"
	case CmdSetField:
		return "Field"
	case CmdSetObjectReference:
		return "ObjectReference"
	case CmdSetStringReference:
		return "StringReference"
	case CmdSetThreadReference:
		return "ThreadReference"
	case CmdSetThreadGroupReference:
		return "ThreadGroupReference"
	case CmdSetArrayReference:
		return "ArrayReference"
	case CmdSetClassLoaderReference:
		return "ClassLoaderReference"
	case CmdSetEventRequest:
		return "EventRequest"
	case CmdSetStackFrame:
		return "StackFrame"
	case CmdSetClassObjectReference:
		return "ClassObjectReference"
	case CmdSetEvent:
		return "Event"
	}
	return fmt.Sprint(int(c))
}

var (
	CmdVirtualMachineVersion               = Command{CmdSetVirtualMachine, 1}
	CmdVirtualMachineClassesBySignature    = Command{CmdSetVirtualMachine, 2}
	CmdVirtualMachineAllClasses            = Command{CmdSetVirtualMachine, 3}
	CmdVirtualMachineAllThreads            = Command{CmdSetVirtualMachine, 4}
	CmdVirtualMachineTopLevelThreadGroups  = Command{CmdSetVirtualMachine, 5}
	CmdVirtualMachineDispose               = Command{CmdSetVirtualMachine, 6}
	CmdVirtualMachineIDSizes               = Command{CmdSetVirtualMachine, 7}
	CmdVirtualMachineSuspend               = Command{CmdSetVirtualMachine, 8}
	CmdVirtualMachineResume                = Command{CmdSetVirtualMachine, 9}
	CmdVirtualMachineExit                  = Command{CmdSetVirtualMachine, 10}
	CmdVirtualMachineCreateString          = Command{CmdSetVirtualMachine, 11}
	CmdVirtualMachineCapabilities          = Command{CmdSetVirtualMachine, 12}
	CmdVirtualMachineClassPaths            = Command{CmdSetVirtualMachine, 13}
	CmdVirtualMachineDisposeObjects        = Command{CmdSetVirtualMachine, 14}
	CmdVirtualMachineHoldEvents            = Command{CmdSetVirtualMachine, 15}
	CmdVirtualMachineReleaseEvents         = Command{CmdSetVirtualMachine, 16}
	CmdVirtualMachineCapabilitiesNew       = Command{CmdSetVirtualMachine, 17}
	CmdVirtualMachineRedefineClasses       = Command{CmdSetVirtualMachine, 18}
	CmdVirtualMachineSetDefaultStratum     = Command{CmdSetVirtualMachine, 19}
	CmdVirtualMachineAllClassesWithGeneric = Command{CmdSetVirtualMachine, 20}

	CmdReferenceTypeSignature            = Command{CmdSetReferenceType, 1}
	CmdReferenceTypeClassLoader          = Command{CmdSetReferenceType, 2}
	CmdReferenceTypeModifiers            = Command{CmdSetReferenceType, 3}
	CmdReferenceTypeFields               = Command{CmdSetReferenceType, 4}
	CmdReferenceTypeMethods              = Command{CmdSetReferenceType, 5}
	CmdReferenceTypeGetValues            = Command{CmdSetReferenceType, 6}
	CmdReferenceTypeSourceFile           = Command{CmdSetReferenceType, 7}
	CmdReferenceTypeNestedTypes          = Command{CmdSetReferenceType, 8}
	CmdReferenceTypeStatus               = Command{CmdSetReferenceType, 9}
	CmdReferenceTypeInterfaces           = Command{CmdSetReferenceType, 10}
	CmdReferenceTypeClassObject          = Command{CmdSetReferenceType, 11}
	CmdReferenceTypeSourceDebugExtension = Command{CmdSetReferenceType, 12}
	CmdReferenceTypeSignatureWithGeneric = Command{CmdSetReferenceType, 13}
	CmdReferenceTypeFieldsWithGeneric    = Command{CmdSetReferenceType, 14}
	CmdReferenceTypeMethodsWithGeneric   = Command{CmdSetReferenceType, 15}

	CmdClassTypeSuperclass   = Command{CmdSetClassType, 1}
	CmdClassTypeSetValues    = Command{CmdSetClassType, 2}
	CmdClassTypeInvokeMethod = Command{CmdSetClassType, 3}
	CmdClassTypeNewInstance  = Command{CmdSetClassType, 4}

	CmdArrayTypeNewInstance = Command{CmdSetArrayType, 1}

	CmdMethodTypeLineTable                = Command{CmdSetMethod, 1}
	CmdMethodTypeVariableTable            = Command{CmdSetMethod, 2}
	CmdMethodTypeBytecodes                = Command{CmdSetMethod, 3}
	CmdMethodTypeIsObsolete               = Command{CmdSetMethod, 4}
	CmdMethodTypeVariableTableWithGeneric = Command{CmdSetMethod, 5}

	CmdObjectReferenceReferenceType     = Command{CmdSetObjectReference, 1}
	CmdObjectReferenceGetValues         = Command{CmdSetObjectReference, 2}
	CmdObjectReferenceSetValues         = Command{CmdSetObjectReference, 3}
	CmdObjectReferenceMonitorInfo       = Command{CmdSetObjectReference, 5}
	CmdObjectReferenceInvokeMethod      = Command{CmdSetObjectReference, 6}
	CmdObjectReferenceDisableCollection = Command{CmdSetObjectReference, 7}
	CmdObjectReferenceEnableCollection  = Command{CmdSetObjectReference, 8}
	CmdObjectReferenceIsCollected       = Command{CmdSetObjectReference, 9}

	CmdStringReferenceValue = Command{CmdSetStringReference, 1}

	CmdThreadReferenceName                    = Command{CmdSetThreadReference, 1}
	CmdThreadReferenceSuspend                 = Command{CmdSetThreadReference, 2}
	CmdThreadReferenceResume                  = Command{CmdSetThreadReference, 3}
	CmdThreadReferenceStatus                  = Command{CmdSetThreadReference, 4}
	CmdThreadReferenceThreadGroup             = Command{CmdSetThreadReference, 5}
	CmdThreadReferenceFrames                  = Command{CmdSetThreadReference, 6}
	CmdThreadReferenceFrameCount              = Command{CmdSetThreadReference, 7}
	CmdThreadReferenceOwnedMonitors           = Command{CmdSetThreadReference, 8}
	CmdThreadReferenceCurrentContendedMonitor = Command{CmdSetThreadReference, 9}
	CmdThreadReferenceStop                    = Command{CmdSetThreadReference, 10}
	CmdThreadReferenceInterrupt               = Command{CmdSetThreadReference, 11}
	CmdThreadReferenceSuspendCount            = Command{CmdSetThreadReference, 12}

	CmdThreadGroupReferenceName     = Command{CmdSetThreadGroupReference, 1}
	CmdThreadGroupReferenceParent   = Command{CmdSetThreadGroupReference, 2}
	CmdThreadGroupReferenceChildren = Command{CmdSetThreadGroupReference, 3}

	CmdArrayReferenceLength    = Command{CmdSetArrayReference, 1}
	CmdArrayReferenceGetValues = Command{CmdSetArrayReference, 2}
	CmdArrayReferenceSetValues = Command{CmdSetArrayReference, 3}

	CmdClassLoaderReferenceVisibleClasses = Command{CmdSetClassLoaderReference, 1}

	CmdEventRequestSet                 = Command{CmdSetEventRequest, 1}
	CmdEventRequestClear               = Command{CmdSetEventRequest, 2}
	CmdEventRequestClearAllBreakpoints = Command{CmdSetEventRequest, 3}

	CmdStackFrameGetValues  = Command{CmdSetStackFrame, 1}
	CmdStackFrameSetValues  = Command{CmdSetStackFrame, 2}
	CmdStackFrameThisObject = Command{CmdSetStackFrame, 3}
	CmdStackFramePopFrames  = Command{CmdSetStackFrame, 4}

	CmdClassObjectReferenceReflectedType = Command{CmdSetClassObjectReference, 1}

	CmdEventComposite = Command{CmdSetEvent, 1}
)

var CmdNames = map[Command]string{}

func init() {
	register := func(c Command, n string) {
		if _, e := CmdNames[c]; e {
			panic("command already registered")
		}
		CmdNames[c] = n
	}
	register(CmdVirtualMachineVersion, "Version")
	register(CmdVirtualMachineClassesBySignature, "ClassesBySignature")
	register(CmdVirtualMachineAllClasses, "AllClasses")
	register(CmdVirtualMachineAllThreads, "AllThreads")
	register(CmdVirtualMachineTopLevelThreadGroups, "TopLevelThreadGroups")
	register(CmdVirtualMachineDispose, "Dispose")
	register(CmdVirtualMachineIDSizes, "IDSizes")
	register(CmdVirtualMachineSuspend, "Suspend")
	register(CmdVirtualMachineResume, "Resume")
	register(CmdVirtualMachineExit, "Exit")
	register(CmdVirtualMachineCreateString, "CreateString")
	register(CmdVirtualMachineCapabilities, "Capabilities")
	register(CmdVirtualMachineClassPaths, "ClassPaths")
	register(CmdVirtualMachineDisposeObjects, "DisposeObjects")
	register(CmdVirtualMachineHoldEvents, "HoldEvents")
	register(CmdVirtualMachineReleaseEvents, "ReleaseEvents")
	register(CmdVirtualMachineCapabilitiesNew, "CapabilitiesNew")
	register(CmdVirtualMachineRedefineClasses, "RedefineClasses")
	register(CmdVirtualMachineSetDefaultStratum, "SetDefaultStratum")
	register(CmdVirtualMachineAllClassesWithGeneric, "AllClassesWithGeneric")

	register(CmdReferenceTypeSignature, "Signature")
	register(CmdReferenceTypeClassLoader, "ClassLoader")
	register(CmdReferenceTypeModifiers, "Modifiers")
	register(CmdReferenceTypeFields, "Fields")
	register(CmdReferenceTypeMethods, "Methods")
	register(CmdReferenceTypeGetValues, "GetValues")
	register(CmdReferenceTypeSourceFile, "SourceFile")
	register(CmdReferenceTypeNestedTypes, "NestedTypes")
	register(CmdReferenceTypeStatus, "Status")
	register(CmdReferenceTypeInterfaces, "Interfaces")
	register(CmdReferenceTypeClassObject, "ClassObject")
	register(CmdReferenceTypeSourceDebugExtension, "SourceDebugExtension")
	register(CmdReferenceTypeSignatureWithGeneric, "SignatureWithGeneric")
	register(CmdReferenceTypeFieldsWithGeneric, "FieldsWithGeneric")
	register(CmdReferenceTypeMethodsWithGeneric, "MethodsWithGeneric")

	register(CmdClassTypeSuperclass, "Superclass")
	register(CmdClassTypeSetValues, "SetValues")
	register(CmdClassTypeInvokeMethod, "InvokeMethod")
	register(CmdClassTypeNewInstance, "NewInstance")

	register(CmdArrayTypeNewInstance, "NewInstance")

	register(CmdMethodTypeLineTable, "LineTable")
	register(CmdMethodTypeVariableTable, "VariableTable")
	register(CmdMethodTypeBytecodes, "Bytecodes")
	register(CmdMethodTypeIsObsolete, "IsObsolete")
	register(CmdMethodTypeVariableTableWithGeneric, "VariableTableWithGeneric")

	register(CmdObjectReferenceReferenceType, "ReferenceType")
	register(CmdObjectReferenceGetValues, "GetValues")
	register(CmdObjectReferenceSetValues, "SetValues")
	register(CmdObjectReferenceMonitorInfo, "MonitorInfo")
	register(CmdObjectReferenceInvokeMethod, "InvokeMethod")
	register(CmdObjectReferenceDisableCollection, "DisableCollection")
	register(CmdObjectReferenceEnableCollection, "EnableCollection")
	register(CmdObjectReferenceIsCollected, "IsCollected")

	register(CmdStringReferenceValue, "Value")

	register(CmdThreadReferenceName, "Name")
	register(CmdThreadReferenceSuspend, "Suspend")
	register(CmdThreadReferenceResume, "Resume")
	register(CmdThreadReferenceStatus, "Status")
	register(CmdThreadReferenceThreadGroup, "ThreadGroup")
	register(CmdThreadReferenceFrames, "Frames")
	register(CmdThreadReferenceFrameCount, "FrameCount")
	register(CmdThreadReferenceOwnedMonitors, "OwnedMonitors")
	register(CmdThreadReferenceCurrentContendedMonitor, "CurrentContendedMonitor")
	register(CmdThreadReferenceStop, "Stop")
	register(CmdThreadReferenceInterrupt, "Interrupt")
	register(CmdThreadReferenceSuspendCount, "SuspendCount")

	register(CmdThreadGroupReferenceName, "Name")
	register(CmdThreadGroupReferenceParent, "Parent")
	register(CmdThreadGroupReferenceChildren, "Children")

	register(CmdArrayReferenceLength, "Length")
	register(CmdArrayReferenceGetValues, "GetValues")
	register(CmdArrayReferenceSetValues, "SetValues")

	register(CmdClassLoaderReferenceVisibleClasses, "VisibleClasses")

	register(CmdEventRequestSet, "Set")
	register(CmdEventRequestClear, "Clear")
	register(CmdEventRequestClearAllBreakpoints, "ClearAllBreakpoints")

	register(CmdStackFrameGetValues, "GetValues")
	register(CmdStackFrameSetValues, "SetValues")
	register(CmdStackFrameThisObject, "ThisObject")
	register(CmdStackFramePopFrames, "PopFrames")

	register(CmdClassObjectReferenceReflectedType, "ReflectedType")

	register(CmdEventComposite, "Composite")
}
