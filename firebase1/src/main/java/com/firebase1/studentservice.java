package com.firebase1;

import com.google.api.core.ApiFuture;
import com.google.cloud.firestore.DocumentReference;
import com.google.cloud.firestore.DocumentSnapshot;
import com.google.cloud.firestore.Firestore;
import com.google.cloud.firestore.WriteResult;
import com.google.firebase.cloud.FirestoreClient;
import org.springframework.stereotype.Service;

import java.util.concurrent.ExecutionException;

@Service
public class studentservice{
public String createCRUD (student student) throws ExecutionException, InterruptedException {
    Firestore dbfirestore= FirestoreClient.getFirestore();
    ApiFuture<WriteResult> collectionApiFuture = dbfirestore.collection("student").document(student.getName()).set(student);
    return collectionApiFuture.get().getUpdateTime().toString();
}
public student getCRUD(String documentid) throws ExecutionException, InterruptedException {
    Firestore dbfirestore= FirestoreClient.getFirestore();
    DocumentReference documentReference = dbfirestore.collection("student").document(documentid);
    ApiFuture<DocumentSnapshot> future= documentReference.get();
    DocumentSnapshot document = future.get();
    student student;
    if(document.exists()) {
        student =document.toObject(student.class);
        return student;
    }
    return null;
}
public String deleteCRUD (String documentid){
    Firestore dbfirestore= FirestoreClient.getFirestore();
    ApiFuture<WriteResult> writeresult= dbfirestore.collection("student").document(documentid).delete();
    return "sucessfully deleted"+documentid;
}
public String updateCRUD (student student){
    Firestore dbfirestore= FirestoreClient.getFirestore();
    ApiFuture<WriteResult> writeresult= dbfirestore.collection("student").document(student.getName()).set(student);
    try {
        return writeresult.get().getUpdateTime().toString();
    } catch (InterruptedException e) {
        throw new RuntimeException(e);
    } catch (ExecutionException e) {
        throw new RuntimeException(e);
    }}}
//public class studentservice {
//    public student getcrud(String documentid) throws ExecutionException, InterruptedException {
//        Firestore firestoredb = FirestoreClient.getFirestore();
//        DocumentReference documentReference = firestoredb.collection("student").document(documentid);
//        ApiFuture<DocumentSnapshot> future= documentReference.get();
//        DocumentSnapshot document = future.get();
//        student student;
//        if(document.exists()) {
//            student =document.toObject(student.class);
//            return student;
//        }
//        return null;
//    }
//
//    public String createcrud(student student) throws ExecutionException, InterruptedException {
//        Firestore firestoredb = FirestoreClient.getFirestore();
//        ApiFuture<WriteResult> writeresult = firestoredb.collection("student").document(student.getDocumentid()).set(student);
//        return writeresult.get().getUpdateTime().toString();
//    }
//
//    public String updatestudent(student student) throws ExecutionException, InterruptedException {
//        Firestore firestoredb = FirestoreClient.getFirestore();
//        ApiFuture<WriteResult> writeresult = firestoredb.collection("student").document(student.getDocumentid()).set(student);
//        return writeresult.get().getUpdateTime().toString();
//    }
//
//    public String deletetudent(String documetid) throws ExecutionException, InterruptedException {
//        Firestore firestoredb = FirestoreClient.getFirestore();
//        ApiFuture<WriteResult> writeresult = firestoredb.collection("student").document(documetid).delete();
//        return writeresult.get().getUpdateTime().toString();
//    }}
