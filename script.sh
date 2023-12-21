i=10
while read -r line ; 
do 
    i=$((i + 1))
    
    echo $line >  testdata/input_tests/test_case_${i}_input.txt ;
done < sample.txt

j=10
while read -r line ; 
do 
    j=$((j + 1))

    
    echo $line >  testdata/output_tests/test_case_${j}_output.txt ;
done < result.txt



    # if [ $j -gt 11 ] ; then
    #     break
    # fi