package Interop

/*
#cgo CFLAGS: -I/usr/lib/jvm/java-11-openjdk-amd64/include -I/usr/lib/jvm/java-11-openjdk-amd64/include/linux
#cgo LDFLAGS: -L/usr/lib/jvm/java-11-openjdk-amd64/lib/amd64/server -ljvm
#include <jni.h> // The Java Native Interface header
#include <stdlib.h>
#include <string.h>

// Global variables
JavaVM* jvm; // a global variable to hold the JVM instance
jclass chordClass;
jmethodID chordConstructorNewChord;
jmethodID chordConstructorJoinChord;
jmethodID methodSet;
jmethodID methodGet;
jmethodID methodDelete;
jmethodID methodGetAllKeys;

char* get_env(JNIEnv** env)
{
	// Get the JVM environment for the current thread
	jint res = (*jvm)->GetEnv(jvm, (void**)env, JNI_VERSION_1_6);
	if (res == JNI_OK) // JVM is attached
		return NULL;

	else if (res == JNI_EDETACHED)
	{
		res = (*jvm)->AttachCurrentThread(jvm, (void **)env, NULL);
		if (res != JNI_OK)
			return strdup("Failed to attach thread to JVM");
		return NULL; // All is well
	}
	else
		return strdup("Failed to get JNIEnv");
}

char* get_exception_message(JNIEnv* env)
{
	// clear exception	
 	jthrowable exception = (*env)->ExceptionOccurred(env);
	(*env)->ExceptionClear(env);
	jclass classThrowable = (*env)->FindClass(env, "java/lang/Throwable");

	// toString()
	jmethodID methodToString = (*env)->GetMethodID(env, classThrowable, "toString", "()Ljava/lang/String;");
	jstring message = (jstring)(*env)->CallObjectMethod(env, exception, methodToString);

	// convert to C string
	const char* messageChars = (*env)->GetStringUTFChars(env, message, NULL);
	char* messageCopy = strdup(messageChars);

	// cleanup
	(*env)->ReleaseStringUTFChars(env, message, messageChars);
	(*env)->DeleteLocalRef(env, message);
	(*env)->DeleteLocalRef(env, classThrowable);
	(*env)->DeleteLocalRef(env, exception);
	return messageCopy;
}

char* init_jvm()
{
	JavaVMInitArgs vm_args; // Initialization arguments vm_args.version = JNI_VERSION_1_6; // set the JNI version vm_args.nOptions = 0; // no options (like class path) vm_args.ignoreUnrecognized = 0;
	JNIEnv* env;

	int res = JNI_CreateJavaVM(&jvm, (void**)&env, &vm_args);
	if (res < 0) 
	{
        char* error_msg;
        switch(res) 
		{
            case JNI_ERR:
                error_msg = "unknown error";
                break;
			case JNI_EDETACHED:
				error_msg = "thread detached from the VM"; break;
            case JNI_EVERSION:
                error_msg = "JNI version error";
                break;
            case JNI_ENOMEM:
                error_msg = "not enough memory";
                break;
            case JNI_EEXIST:
                error_msg = "VM already created";
                break;
            case JNI_EINVAL:
                error_msg = "invalid arguments";
                break;
            default:
                error_msg = "unknown error code";
		}
        return error_msg;
    }
	
    return NULL;
}

char* load_chord_class()
{
	JNIEnv* env = NULL;
	chordClass = NULL;
	chordConstructorNewChord = NULL;
	chordConstructorJoinChord = NULL;
	methodSet = NULL;
	methodGet = NULL;
	methodDelete = NULL;
	methodGetAllKeys = NULL;

	char* error = get_env(&env);
	if (error != NULL)
		return error;

	// [Java] URL url = new URL("file:///workspaces/large-scale-workshop/interop/");

	jclass urlClass = (*env)->FindClass(env, "java/net/URL");
	if (urlClass == NULL)
	{
		error = "Could not find URL class";
		goto cleanup;
	}
	
	jmethodID urlConstructor = (*env)->GetMethodID(env, urlClass, "<init>", "(Ljava/lang/String;)V");
	if (urlConstructor == NULL)
	{
		error = "Could not find URL constructor";
		goto cleanup;
	}

	jstring urlStr = (*env)->NewStringUTF(env, "file:///workspaces/large-scale-workshop/interop/");
	jobject urlObj = (*env)->NewObject(env, urlClass, urlConstructor, urlStr);
	if (urlObj == NULL)
	{
		error = "Could not create instance of URL object";
		goto cleanup;
	}

	// [Java] URLClassLoader urlClassLoader = new URLClassLoader(new URL[]{url});

	jobjectArray urlArray = (*env)->NewObjectArray(env, 1, urlClass, urlObj); if (urlObj == NULL)
	{
		error = "Could not create instance of URL[]{urlObj}";
    	goto cleanup;
	}

	jclass urlClassLoaderClass = (*env)->FindClass(env, "java/net/URLClassLoader");
	if (urlClassLoaderClass == NULL)
	{
		error = "Could not find URLClassLoader class";
		goto cleanup;
	}

	jmethodID urlClassLoaderConstructor = (*env)->GetMethodID(env, urlClassLoaderClass, "<init>", "([Ljava/net/URL;)V");
	if (urlClassLoaderConstructor == NULL)
	{
		error = "Could not find URLClassLoader constructor";
		goto cleanup;
	}

	// [Java] Class<?> chordClass = urlClassLoader.loadClass("dht.Chord");

	jobject urlClassLoader = (*env)->NewObject(env, urlClassLoaderClass, urlClassLoaderConstructor, urlArray);
	jmethodID loadClassMethod = (*env)->GetMethodID(env, urlClassLoaderClass, "loadClass", "(Ljava/lang/String;)Ljava/lang/Class;");
	if (loadClassMethod == NULL)
	{
		error = "Could not find loadClass method";
		goto cleanup;
	}

	jstring className = (*env)->NewStringUTF(env, "dht.Chord");
	jclass localChordClass = (jclass)(*env)->CallObjectMethod(env, urlClassLoader, loadClassMethod, className);
	if (localChordClass == NULL)
	{
		if ((*env)->ExceptionCheck(env))
			error = get_exception_message(env);
		else
			error = "Could not load Chord class";
		goto cleanup;
	}

	// set the global reference
	chordClass = (*env)->NewGlobalRef(env, localChordClass);


	// load the object methods so we can call them from Go
	chordConstructorNewChord = (*env)->GetMethodID(env, chordClass, "<init>", "(Ljava/lang/String;I)V");
	if ((*env)->ExceptionCheck(env))
	{
		error = get_exception_message(env);
		goto cleanup;
	}
	chordConstructorJoinChord = (*env)->GetMethodID(env, chordClass, "<init>", "(Ljava/lang/String;Ljava/lang/String;I)V");
	if ((*env)->ExceptionCheck(env))
	{
		error = get_exception_message(env);
		goto cleanup;
	}
	methodSet = (*env)->GetMethodID(env, chordClass, "set", "(Ljava/lang/String;Ljava/lang/String;)V");
	if ((*env)->ExceptionCheck(env))
	{
		error = get_exception_message(env);
		goto cleanup;
	}
	methodGet = (*env)->GetMethodID(env, chordClass, "get", "(Ljava/lang/String;)Ljava/lang/String;");
	if ((*env)->ExceptionCheck(env))
	{
		error = get_exception_message(env);
		goto cleanup;
	}
	methodDelete = (*env)->GetMethodID(env, chordClass, "delete", "(Ljava/lang/String;)V");
	if ((*env)->ExceptionCheck(env))
	{
		error = get_exception_message(env);
		goto cleanup;
	}
	methodGetAllKeys = (*env)->GetMethodID(env, chordClass, "getAllKeys", "()[Ljava/lang/String;");
	if ((*env)->ExceptionCheck(env))
	{
		error = get_exception_message(env);
		goto cleanup;
	}

cleanup:
	if (urlClass != NULL) 
		(*env)->DeleteLocalRef(env, urlClass);
	if (urlObj != NULL) 
		(*env)->DeleteLocalRef(env, urlObj);
	if (urlArray != NULL) 
		(*env)->DeleteLocalRef(env, urlArray);
	if (urlClassLoaderClass != NULL) 
		(*env)->DeleteLocalRef(env, urlClassLoaderClass); 
	if (urlClassLoader != NULL) 
		(*env)->DeleteLocalRef(env, urlClassLoader);
	if (className != NULL) 
		(*env)->DeleteLocalRef(env, className);
	if (localChordClass != NULL) 
		(*env)->DeleteLocalRef(env, localChordClass);
	return error;
}

// Java API: public Chord(String node_name, int chord_port);
jobject call_chord_constructor_new_chord(char* node_name, int port, char** out_error)
{
    jobject chordObject = NULL;
    jstring jnodeName = NULL;
    JNIEnv* env;

    char* error = get_env(&env);
    if (error != NULL)
    {
		*out_error = error;
        return NULL;
    }

	// Convert the C string to a Java string
	jnodeName = (*env)->NewStringUTF(env, node_name);
	if ((*env)->ExceptionCheck(env))
	{
		*out_error = get_exception_message(env);
		goto cleanup;
	}

	// Call the constructor using NewObject function
	chordObject = (*env)->NewObject(env, chordClass, chordConstructorNewChord, jnodeName, port); if ((*env)->ExceptionCheck(env))
	{
		*out_error = get_exception_message(env);
		goto cleanup;
	}

cleanup:
	if (jnodeName)
		(*env)->DeleteLocalRef(env, jnodeName);
	return chordObject;
}

// Java API: public Chord(String node_name, String root_name, int chord_port);
jobject call_chord_constructor_join_chord(char* node_name, char* root_name, int port, char** out_error)
{
	jobject newChordObject = NULL;
	jstring jnodeName = NULL;
	jstring jrootName = NULL;
	JNIEnv* env = NULL;

	char* error = get_env(&env);
	if (error != NULL)
	{
		*out_error = error;
		return NULL;
	}

	// Convert the C strings to Java strings
	jnodeName = (*env)->NewStringUTF(env, node_name);
	jrootName = (*env)->NewStringUTF(env, root_name);
	if ((*env)->ExceptionCheck(env))
	{
		*out_error = get_exception_message(env);
		goto cleanup;
	}

	// Call the constructor
	newChordObject = (*env)->NewObject(env, chordClass, chordConstructorJoinChord, jnodeName, jrootName, port);
	if ((*env)->ExceptionCheck(env))
	{
		*out_error = get_exception_message(env);
		goto cleanup;
	}

cleanup:
	if (jnodeName)
		(*env)->DeleteLocalRef(env, jnodeName);
	if (jrootName)
		(*env)->DeleteLocalRef(env, jrootName);
	return newChordObject;
}

// Java API: public void set(String key, String val);
void call_method_set(jobject chordObject, char* key, char* value, char** out_error)
{
    jstring jkey = NULL;
    jstring jvalue = NULL;
	JNIEnv* env;

	char* error = get_env(&env);
	if (error != NULL)
	{
        *out_error = error;
		return;
	}

    // Convert the C strings to Java strings
    jkey = (*env)->NewStringUTF(env, key);
    jvalue = (*env)->NewStringUTF(env, value);
    if ((*env)->ExceptionCheck(env))
	{
		*out_error = get_exception_message(env);
        goto cleanup;
    }

	(*env)->CallVoidMethod(env, chordObject, methodSet, jkey, jvalue);
	if ((*env)->ExceptionCheck(env))
		*out_error = get_exception_message(env);

cleanup:
    if (jkey)
		(*env)->DeleteLocalRef(env, jkey);
    if (jvalue)
		(*env)->DeleteLocalRef(env, jvalue);
}

// Java API: public String get(String key);
char* call_method_get(jobject chordObject, char* key, char** out_error)
{
	jstring jkey = NULL; 
	jstring jresult = NULL; 
	const char* result = NULL;
	JNIEnv* env;

	char* error = get_env(&env);
	if (error != NULL)
	{
        *out_error = error;
		return;
	}
	
	// Convert the C string to a Java string
	jkey = (*env)->NewStringUTF(env, key);
	if ((*env)->ExceptionCheck(env))
	{
		*out_error = get_exception_message(env);
		goto cleanup;
	}

	jresult = (jstring)(*env)->CallObjectMethod(env, chordObject, methodGet, jkey);
	if ((*env)->ExceptionCheck(env))
	{
		*out_error = get_exception_message(env);
		goto cleanup;
	}
	if (jresult == NULL) // not found
		goto cleanup;

	// Convert the Java string to a C string
	result = (*env)->GetStringUTFChars(env, jresult, 0);
	if ((*env)->ExceptionCheck(env))
		*out_error = get_exception_message(env);

cleanup:
	if (jkey)
		(*env)->DeleteLocalRef(env, jkey);
	if (jresult)
		(*env)->DeleteLocalRef(env, jresult);
	return (char*)result;
}

// Java API: public void delete(String key);
void call_method_delete(jobject chordObject, char* key, char** out_error)
{
    jstring jkey = NULL;
	JNIEnv* env;

	char* error = get_env(&env);
	if (error != NULL)
	{
        *out_error = error;
		return;
	}

    // Convert the C string to Java string
    jkey = (*env)->NewStringUTF(env, key);
    if ((*env)->ExceptionCheck(env))
	{
		*out_error = get_exception_message(env);
        goto cleanup;
    }

	(*env)->CallVoidMethod(env, chordObject, methodDelete, jkey);
	if ((*env)->ExceptionCheck(env))
		*out_error = get_exception_message(env);

cleanup:
    if (jkey)
		(*env)->DeleteLocalRef(env, jkey);
}

// Java API: public String[] getAllKeys();
char** call_method_get_all_keys(jobject chordObject, char** out_error)
{
	jobjectArray jresult = NULL;
	char** result = NULL;
	JNIEnv* env;

	char* error = get_env(&env);
	if (error != NULL)
	{
		*out_error = error;
		return NULL;
	}

	jsize len;
	int i;

	jresult = (jobjectArray)(*env)->CallObjectMethod(env, chordObject, methodGetAllKeys);
	if ((*env)->ExceptionCheck(env))
	{
		*out_error = get_exception_message(env);
		goto cleanup;
	}

	// Convert the Java string array to a C string array
	len = (*env)->GetArrayLength(env, jresult);

	// copy the jresult to “result”


	// TODO: Make sure you allocate enough space!
	// TODO: Make sure the LAST element in the array is NULL – to mark it is the end of the array


cleanup:
	if (jresult) 
		(*env)->DeleteLocalRef(env, jresult);
	return result;
}
*/
import "C"