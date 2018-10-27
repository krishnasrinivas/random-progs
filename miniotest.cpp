#include "aws/core/Aws.h"
#include "aws/core/auth/AWSCredentialsProvider.h"
#include "aws/core/client/ClientConfiguration.h"
#include "aws/s3/S3Client.h"
#include "aws/s3/model/CreateBucketRequest.h"
#include "aws/s3/model/PutObjectRequest.h"
#include <fstream>
#include <iostream>

using namespace Aws;

int main()
{
	//Init SDK
	SDKOptions options;
	options.loggingOptions.logLevel = Utils::Logging::LogLevel::Trace;
	InitAPI(options);

	//Init Client
	Client::ClientConfiguration config;
	config.region = "us-east-1";
	// config.scheme = Http::Scheme::HTTP;
	config.endpointOverride = "play.minio.io:9000";
	// config.verifySSL = false;
	S3::S3Client client(Aws::Auth::AWSCredentials("Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"), config, Aws::Client::AWSAuthV4Signer::PayloadSigningPolicy::Never, false);

	//Create Bucket
	String mybucket = "awstominio";
	S3::Model::CreateBucketRequest create_request;
	create_request.SetBucket(mybucket);
	std::cout<<"Creating: "<<mybucket << "\n";
	auto result = client.CreateBucket(create_request);
	if(result.IsSuccess())
		std::cout<<"CREATED: "<<mybucket << "\n";
	else
		std::cout<<"Error: "<<result.GetError().GetMessage() <<"\n";

	//List Buckets
/*	auto listresult = client.ListBuckets();
	if(listresult.IsSuccess())
	{
		Vector<S3::Model::Bucket> buckets = listresult.GetResult().GetBuckets();
		for(auto const &bucket:buckets)
			std::cout<<bucket.GetName()<<"\n";
	}
	else
	{
		std::cout<<"List error: "<<listresult.GetError().GetMessage() <<"\n";
	}
*/
	//Put Object
	String mykey = "myFilePath";
	String myobject = "/etc/passwd";
	S3::Model::PutObjectRequest put_req;
	put_req.WithBucket(mybucket).WithKey("somekey1");
	// put_req.SetBucket(mybucket);
	// put_req.SetKey(mykey);
	// put_req.SetContentType("application/txt");
	auto input = MakeShared<FStream>("PutObjectInputStream",myobject.c_str(),std::ios_base::in|std::ios::binary);
	// char cline[20];
	// input->getline(cline,19);
	// std::cout<<cline<<std::endl;
	put_req.SetBody(input);
	auto put_res = client.PutObject(put_req);
	if(put_res.IsSuccess())
		std::cout<<"PUT SUCCESS!\n";
	else
	{
		std::cout<<"Put Error: "<<put_res.GetError().GetMessage() <<"\n";
		std::cout<<"Put Error: "<<put_res.GetError().GetExceptionName() <<"\n";
	}

	ShutdownAPI(options);
	return 0;
}
