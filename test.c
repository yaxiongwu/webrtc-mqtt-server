#include <stdio.h>
#include <string.h>
int main(){
   char *p1="fromClinet/c12456ewq/answer";
   const char* p2 = "/";
   char s[3][12],s1[12],s2[12],s3[12];
	char* ret1 = strstr(p1, p2);//把返回的字符串首地址赋给ret
	if (ret1 == NULL)
	{
		printf("子串不存在\n");//当返回的字符串首地址为空，ret为一个空指针，代表不存在该子串
	}
	else
	{
		printf("%s,ret1-p1:%d\n", ret1+1,ret1-p1);//当返回的字符串首地址不为空，则会从字符串首地址开始打印，到‘\0’停止
        snprintf(s1,ret1-p1+1,p1);
	}
    
    char* ret2 =strstr(ret1+1,p2);
    if (ret2 == NULL)
	{
		printf("子串不存在\n");//当返回的字符串首地址为空，ret为一个空指针，代表不存在该子串
	}
	else
	{
		printf("%s,ret2-ret1:%d\n", ret2+1,ret2-ret1);//当返回的字符串首地址不为空，则会从字符串首地址开始打印，到‘\0’停止
        snprintf(s2,ret2-ret1,ret1+1);
        strcpy(s3,ret2+1);
	}

   
    printf(" s1:%s \n s2:%s \n s3:%s\n", s1,s2,s3);//当返回的字符串首地址不为空，则会从字符串首地址开始打印，到‘\0’停止
    return 0;
}